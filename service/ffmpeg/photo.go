package service

import (
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/core/repository"
	"github.com/41x3n/Xom/core/usecase"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pdfcpuAPI "github.com/pdfcpu/pdfcpu/pkg/api"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
	"time"
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

	if pu.UpdatePhotoStatus(photo, domain.Processing) != nil {
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
		errConvert = c.ConvertPhoto(inputPath, outputPath)
	}
	if errConvert != nil {
		return errConvert
	}

	errorSendPhoto := c.SendFileToUser(photo, outputPath,
		"Here is your converted photo")
	if errorSendPhoto != nil {
		return errorSendPhoto
	}

	if pu.UpdatePhotoStatus(photo, domain.Completed) != nil {
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

func (c *converter) ConvertPhoto(inputPath, outputPath string) error {
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

func (c *converter) SendFileToUser(photo *domain.Photo, outputPath, message string) error {
	telegramAPI := c.telegram.GetAPI()

	userTelegramID := photo.UserTelegramID
	fileFormat := photo.ConvertTo

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

	// Send the document
	_, err := telegramAPI.Send(documentToBeSent)
	if err != nil {
		return err
	}

	return nil
}
