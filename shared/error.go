package shared

import (
	"errors"
	"fmt"
)

var ErrInvalidFile = errors.New("invalid file")
var ErrInvalidCallbackData = errors.New("invalid callback data")
var ErrUpdateStatus = errors.New("error updating status")
var ErrFileIsAlreadyProcessing = errors.New("file is already processing")
var ErrFileAlreadyProcessed = errors.New("file already processed")
var ErrFileFailed = errors.New("file processing failed")
var ErrFileFormatInvalid = errors.New("file format is invalid")
var ErrFailedToConvert = errors.New("failed to convert file, " +
	"please try again later")

var StatusToError = map[Status]error{
	Processing: ErrFileIsAlreadyProcessing,
	Completed:  ErrFileAlreadyProcessed,
	Failed:     ErrFileFailed,
}

func HandleFileStateError(status Status) error {
	if err, ok := StatusToError[status]; ok {
		return err
	}
	return fmt.Errorf("unhandled status: %v", status)
}
