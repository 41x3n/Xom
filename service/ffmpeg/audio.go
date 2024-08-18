package service

import (
	"log"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/core/repository"
	"github.com/41x3n/Xom/core/usecase"
	"github.com/41x3n/Xom/shared"
)

func (c *converter) HandleAudios(ID int64) error {
	contextTimeout := time.Duration(c.env.ContextTimeout) * time.Second

	ar := repository.NewAudioRepository(c.db, domain.TableAudio)
	au := usecase.NewAudioUsecase(ar, contextTimeout)

	log.Printf("Handling audios for ID: %d", ID)

	audio, err := au.ValidateIfAudioReadyToBeConverted(ID)
	if err != nil {
		return err
	}

	errFormatValid := c.IsValidAudioFormat(audio.FileType)
	if !errFormatValid {
		return shared.ErrFileFormatInvalid
	}

	if au.UpdateAudioStatus(audio, shared.Processing) != nil {
		return shared.ErrUpdateStatus
	}

	inputPath, outputPath, errPath := c.GetInputOutputFilePaths(audio)
	if errPath != nil {
		return errPath
	}

	log.Printf("File Path: %s - %s", inputPath, outputPath)

	errConvert := c.ConvertFile(inputPath, outputPath)
	if errConvert != nil {
		var result *multierror.Error

		errUpdateStatus := au.UpdateAudioStatus(audio, shared.Failed)
		if errUpdateStatus != nil {
			result = multierror.Append(result, errUpdateStatus)
		}

		errInformUser := c.InformUserAboutError(audio.UserTelegramID, audio.MessageID,
			"Sorry, "+shared.ErrFailedToConvert.Error())
		if errInformUser != nil {
			result = multierror.Append(result, errInformUser)
		}

		// Append errConvert to the result
		result = multierror.Append(result, errConvert)

		return result.ErrorOrNil()
	}

	errorSendAudio := c.SendFileToUser(audio, outputPath,
		"✨ Here’s your converted audio!")
	if errorSendAudio != nil {
		return errorSendAudio
	}

	if au.UpdateAudioStatus(audio, shared.Completed) != nil {
		return shared.ErrUpdateStatus
	}

	return nil
}
