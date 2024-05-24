package controller

import (
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackController struct {
	cu          domain.CallbackUseCase
	TelegramAPI *tgbotapi.BotAPI
	RabbitMQ    shared.RabbitMQService
}

func NewCallbackController(cu domain.CallbackUseCase,
	telegramAPI *tgbotapi.BotAPI, rabbitMQ shared.RabbitMQService) *CallbackController {
	return &CallbackController{
		cu:          cu,
		TelegramAPI: telegramAPI,
		RabbitMQ:    rabbitMQ}
}

func (cc *CallbackController) HandleCallback(callback *tgbotapi.CallbackQuery) error {
	photoId, command, convertTo, err := cc.cu.GetFileIDAndCommand(callback)
	if err != nil {
		return err
	}

	if command == string(shared.PhotoCommand) {
		photo, err := cc.cu.GetPhotoByID(photoId)
		if err != nil {
			return err
		}

		if photo.Status == domain.Processing {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID,
				"Your photo is already being processed")
			if _, err = cc.TelegramAPI.Send(msg); err != nil {
				return err
			}
			return nil
		}

		payload := shared.RabbitMQPayload{
			Command: shared.PhotoCommand,
			ID:      photo.ID,
		}

		photo.Status = domain.Preparing
		photo.ConvertTo = convertTo
		if err = cc.cu.MarkPhotoAsPreparing(photo); err != nil {
			return err
		}

		if err = cc.RabbitMQ.PublishMessage(payload); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Your photo is being processed")
		if _, err = cc.TelegramAPI.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
