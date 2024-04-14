package controller

import (
	"bytes"
	"fmt"
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
)

type PhotoController struct {
	TelegramAPI *tgbotapi.BotAPI
}

func NewPhotoController(telegramAPI *tgbotapi.BotAPI) *PhotoController {
	return &PhotoController{TelegramAPI: telegramAPI}
}

func (pc *PhotoController) HandlePhotoCommand(user *domain.User, message *tgbotapi.Message) error {
	var photoID string

	if message.Photo != nil && len(message.Photo) > 0 {
		photoID = message.Photo[3].FileID
	}

	if message.Document != nil {
		photoID = message.Document.FileID
	}

	fmt.Println(photoID)

	fileLink, err := pc.TelegramAPI.GetFileDirectURL(photoID)

	if err != nil {
		return err
	}

	fmt.Println(fileLink)

	// Send a GET request to the file URL
	resp, respErr := http.Get(fileLink)
	if respErr != nil {
		return respErr
	}

	defer shared.CloseResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: Unexpected status code: %d", resp.StatusCode)
	}

	// Create a buffer to store the image data
	var imageBuffer bytes.Buffer

	// Copy the response body into the buffer
	_, err = io.Copy(&imageBuffer, resp.Body)
	if err != nil {
		return fmt.Errorf("error: Unable to copy image data: %v", err)
	}

	// Print the size of the byte slice
	fmt.Printf("Image size: %d bytes\n", imageBuffer.Len())

	photoConfig := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileBytes{
		Name:  "image.jpg",
		Bytes: imageBuffer.Bytes(),
	})

	_, err = pc.TelegramAPI.Send(photoConfig)

	return err
}
