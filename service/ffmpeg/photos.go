package service

import (
	"github.com/41x3n/Xom/core/domain"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
)

func HandlePhotos(ID int64) error {
	log.Printf("Handling photos for ID: %d", ID)

	return nil
}

func convertFile(inputFile string, inputFormat string, outputFile string,
	outputFormat string) {
	if !isValidFormat(inputFormat) || !isValidFormat(outputFormat) {
		log.Fatalf("Invalid format. Supported formats are jpg, jpeg, png, gif, pdf, webp")
		return
	}

	err := ffmpeg.
		Input(inputFile).
		Output(outputFile, ffmpeg.KwArgs{"f": outputFormat}).
		Run()

	if err != nil {
		log.Fatalf("Error while converting %s to %s: %v", inputFormat, outputFormat, err)
	}
}

func isValidFormat(format string) bool {
	for _, validFormat := range domain.FileTypeArray {
		if format == validFormat {
			return true
		}
	}
	return false
}
