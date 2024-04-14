package controller

import (
	"github.com/41x3n/Xom/core/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PhotoController struct {
	pu          domain.PhotoUseCase
	TelegramAPI *tgbotapi.BotAPI
}

func NewPhotoController(pu domain.PhotoUseCase, telegramAPI *tgbotapi.BotAPI) *PhotoController {
	return &PhotoController{
		pu:          pu,
		TelegramAPI: telegramAPI}
}

func (pc *PhotoController) HandlePhotoCommand(user *domain.User, message *tgbotapi.Message) error {
	fileID, fileType, err := pc.pu.GetFileIDAndType(message)
	if err != nil {
		return err
	}

	_, err = pc.pu.SavePhotoId(user, fileID, fileType)
	if err != nil {
		return err
	}

	buttonRows := pc.pu.GenerateConvertOptions(fileType)
	if len(buttonRows) == 0 {
		return domain.ErrInvalidFile
	}

	keyboard := pc.pu.GenerateKeyboardMarkup(buttonRows)
	msg := pc.pu.GenerateMessage(fileType, message, keyboard)

	_, err = pc.TelegramAPI.Send(msg)
	if err != nil {
		return err
	}

	//fileLink, err := pc.TelegramAPI.GetFileDirectURL(FileID)
	//
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Println(fileLink)
	//
	//// Send a GET request to the file URL
	//resp, respErr := http.Get(fileLink)
	//if respErr != nil {
	//	return respErr
	//}
	//
	//defer shared.CloseResponseBody(resp.Body)
	//
	//if resp.StatusCode != http.StatusOK {
	//	return fmt.Errorf("error: Unexpected status code: %d", resp.StatusCode)
	//}
	//
	//// Create a buffer to store the image data
	//var imageBuffer bytes.Buffer
	//
	//// Copy the response body into the buffer
	//_, err = io.Copy(&imageBuffer, resp.Body)
	//if err != nil {
	//	return fmt.Errorf("error: Unable to copy image data: %v", err)
	//}
	//
	//// Print the size of the byte slice
	//fmt.Printf("Image size: %d bytes\n", imageBuffer.Len())
	//
	//photoConfig := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileBytes{
	//	Name:  "image.jpg",
	//	Bytes: imageBuffer.Bytes(),
	//})
	//
	//_, err = pc.TelegramAPI.Send(photoConfig)

	return nil
}
