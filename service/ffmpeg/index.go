package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
	"gorm.io/gorm"
)

type converter struct {
	env      *shared.Env
	telegram shared.TelegramService
	db       *gorm.DB
}

func (c *converter) HandleFiles(payload *shared.RabbitMQPayload) error {
	var err error

	switch payload.Command {
	case shared.PhotoCommand:
		err = c.HandlePhotos(payload.ID)
	default:
		log.Printf("Unknown command: %s", payload.Command)
	}

	return err
}

func (c *converter) GetInputOutputFilePaths(photo *domain.Photo) (string, string, error) {
	telegramAPI := c.telegram.GetAPI()

	fileID := photo.FileID
	id := photo.ID
	fileType := photo.FileType
	convertTo := photo.ConvertTo

	fileLink, err := telegramAPI.GetFileDirectURL(fileID)
	if err != nil {
		return "", "", err
	}

	resp, respErr := http.Get(fileLink)
	if respErr != nil {
		return "", "", respErr
	}

	defer shared.CloseResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("error: Unexpected status code: %d",
			resp.StatusCode)
	}

	// Check if download folder exists. If not, create it.
	if _, err := os.Stat(shared.DownloadFolder); os.IsNotExist(err) {
		err := os.Mkdir(shared.DownloadFolder, 0755)
		if err != nil {
			return "", "", err
		}
	}

	// Save the file in the download folder
	inputPath := filepath.Join(shared.DownloadFolder, strconv.FormatInt(id, 10)+"."+fileType)
	file, errFilePath := os.Create(inputPath)
	if errFilePath != nil {
		return "", "", err
	}
	defer shared.CloseFile(file)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", "", err
	}

	// Check if converted folder exists. If not, create it.
	if _, err := os.Stat(shared.ConvertedFolder); os.IsNotExist(err) {
		err := os.Mkdir(shared.ConvertedFolder, 0755)
		if err != nil {
			return "", "", err
		}
	}

	outputPath := filepath.Join(shared.ConvertedFolder, strconv.FormatInt(id, 10)+"."+convertTo)

	// Return the path to the file
	return inputPath, outputPath, nil
}

func NewFFMPEGService(env *shared.Env, telegram shared.TelegramService, db *gorm.DB) shared.FFMPEGService {
	return &converter{
		env:      env,
		telegram: telegram,
		db:       db,
	}
}
