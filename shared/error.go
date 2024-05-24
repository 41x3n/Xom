package shared

import (
	"errors"
	"fmt"
	"github.com/41x3n/Xom/core/domain"
)

var ErrInvalidFile = errors.New("invalid file")
var ErrInvalidCallbackData = errors.New("invalid callback data")
var ErrUpdateStatus = errors.New("error updating status")
var ErrFileIsAlreadyProcessing = errors.New("file is already processing")
var ErrFileAlreadyProcessed = errors.New("file already processed")
var ErrFileFailed = errors.New("file processing failed")
var ErrFileFormatInvalid = errors.New("file format is invalid")

var StatusToError = map[domain.Status]error{
	domain.Processing: ErrFileIsAlreadyProcessing,
	domain.Completed:  ErrFileAlreadyProcessed,
	domain.Failed:     ErrFileFailed,
}

func HandleFileStateError(status domain.Status) error {
	if err, ok := StatusToError[status]; ok {
		return err
	}
	return fmt.Errorf("unhandled status: %v", status)
}
