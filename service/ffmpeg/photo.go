package service

import (
	"log"
	"os"
	"time"

	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/core/repository"
	"github.com/41x3n/Xom/core/usecase"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hashicorp/go-multierror"
	pdfcpuAPI "github.com/pdfcpu/pdfcpu/pkg/api"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func (c *converter) HandlePhotos(ID int64) error {
	contextTimeout := time.Duration(c.env.ContextTimeout) * time.Second

	pr := repository.NewPhotoRepository(c.db, domain.TablePhoto)
	pu := usecase.NewPhotoUsecase(pr, contextTimeout)

	log.Printf("Handling photos for ID: %d", ID)

	photo, err := pu.ValidateIfPhotoReadyToBeConverted(ID)
	if err != nil {
		return err
	}

	errFormatValid := c.IsValidFormat(photo.FileType)
	if !errFormatValid {
		return shared.ErrFileFormatInvalid
	}

	if pu.UpdatePhotoStatus(photo, shared.Processing) != nil {
		return shared.ErrUpdateStatus
	}

	inputPath, outputPath, errPath := c.GetInputOutputFilePaths(photo)
	if errPath != nil {
		return errPath
	}

	log.Printf("File Path: %s - %s", inputPath, outputPath)

	var errConvert error
	if photo.ConvertTo == "pdf" {
		errConvert = c.ConvertImageToPDF(inputPath, outputPath)
	} else {
		errConvert = c.ConvertFile(inputPath, outputPath)
	}
	if errConvert != nil {
		var result *multierror.Error

		errUpdateStatus := pu.UpdatePhotoStatus(photo, shared.Failed)
		if errUpdateStatus != nil {
			result = multierror.Append(result, errUpdateStatus)
		}

		errInformUser := c.InformUserAboutError(photo.UserTelegramID,
			photo.MessageID,
			"Sorry, "+shared.ErrFailedToConvert.Error())
		if errInformUser != nil {
			result = multierror.Append(result, errInformUser)
		}

		// Append errConvert to the result
		result = multierror.Append(result, errConvert)

		return result.ErrorOrNil()
	}

	errorSendPhoto := c.SendFileToUser(photo, outputPath,
		"Here is your converted photo")
	if errorSendPhoto != nil {
		return errorSendPhoto
	}

	if pu.UpdatePhotoStatus(photo, shared.Completed) != nil {
		return shared.ErrUpdateStatus
	}

	return nil
}

func (c *converter) IsValidFormat(format string) bool {
	for _, validFormat := range domain.FileTypeArray {
		if format == validFormat {
			return true
		}
	}
	return false
}

func (c *converter) IsValidAudioFormat(format string) bool {
	for _, validFormat := range domain.AudioFileTypeArray {
		if format == validFormat {
			return true
		}
	}
	return false
}

func (c *converter) ConvertFile(inputPath, outputPath string) error {
	err := ffmpeg.Input(inputPath).Output(outputPath).OverWriteOutput().Run()

	return err
}

func (c *converter) ConvertImageToPDF(inputPath, outputPath string) error {
	errPDF := pdfcpuAPI.ImportImagesFile([]string{inputPath}, outputPath, nil,
		nil)

	time.Sleep(3 * time.Second)
	if errPDF != nil {
		return errPDF
	}

	return nil
}

func (c *converter) SendFileToUser(media interface{}, outputPath, message string) error {
	telegramAPI := c.telegram.GetAPI()

	var userTelegramID, messageId int64
	var fileFormat string

	switch v := media.(type) {
	case *domain.Photo:
		userTelegramID = v.UserTelegramID
		fileFormat = v.ConvertTo
		messageId = v.MessageID
	case *domain.Audio:
		userTelegramID = v.UserTelegramID
		fileFormat = v.ConvertTo
		messageId = v.MessageID
	}

	file, errFileOpen := os.Open(outputPath)
	if errFileOpen != nil {
		return errFileOpen
	}

	defer shared.CloseFile(file)

	reader := tgbotapi.FileReader{
		Name:   "converted." + fileFormat,
		Reader: file,
	}
	documentToBeSent := tgbotapi.NewDocument(userTelegramID, reader)
	documentToBeSent.Caption = message
	documentToBeSent.ReplyToMessageID = int(messageId)

	// Send the document
	_, err := telegramAPI.Send(documentToBeSent)
	if err != nil {
		return err
	}

	return nil
}

func (c *converter) InformUserAboutError(userTelegramID, messageID int64,
	errorText string) error {
	telegramAPI := c.telegram.GetAPI()

	msg := tgbotapi.NewMessage(userTelegramID, errorText)
	msg.ReplyToMessageID = int(messageID)
	_, err := telegramAPI.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
