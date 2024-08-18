package controller

import (
	"github.com/41x3n/Xom/core/domain"
	interfaces "github.com/41x3n/Xom/interface"
	"github.com/41x3n/Xom/shared"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackController struct {
	cu          domain.CallbackUseCase
	TelegramAPI *tgbotapi.BotAPI
	RabbitMQ    interfaces.RabbitMQService
}

func NewCallbackController(cu domain.CallbackUseCase,
	telegramAPI *tgbotapi.BotAPI, rabbitMQ interfaces.RabbitMQService) *CallbackController {
	return &CallbackController{
		cu:          cu,
		TelegramAPI: telegramAPI,
		RabbitMQ:    rabbitMQ}
}

func (cc *CallbackController) HandleCallback(callback *tgbotapi.CallbackQuery) error {
	fileId, command, convertTo, err := cc.cu.GetFileIDAndCommand(callback)
	if err != nil {
		return err
	}

	if command == string(shared.PhotoCommand) {
		photo, err := cc.cu.GetPhotoByID(fileId)
		if err != nil {
			return err
		}

		if photo.Status == shared.Processing {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID,
				"Hey, your photo is still being processed. Please wait a bit longer.")
			msg.ReplyToMessageID = int(photo.MessageID)
			if _, err = cc.TelegramAPI.Send(msg); err != nil {
				return err
			}
			return nil
		}

		payload := shared.RabbitMQPayload{
			Command: shared.PhotoCommand,
			ID:      photo.ID,
		}

		photo.Status = shared.Preparing
		photo.ConvertTo = convertTo
		if err = cc.cu.MarkPhotoAsPreparing(photo); err != nil {
			return err
		}

		if err = cc.RabbitMQ.PublishMessage(payload); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "⏳ Hold tight..., "+
			"your photo is being processed...")
		msg.ReplyToMessageID = int(photo.MessageID)
		if _, err = cc.TelegramAPI.Send(msg); err != nil {
			return err
		}
	}

	if command == string(shared.AudioCommand) {
		audio, errOnFind := cc.cu.GetAudioByID(fileId)
		if errOnFind != nil {
			return errOnFind
		}

		if audio.Status == shared.Processing {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID,
				"Hey, your audio is still being processed. Please wait a bit longer.")
			msg.ReplyToMessageID = int(audio.MessageID)
			if _, errOnResponse := cc.TelegramAPI.Send(msg); errOnResponse != nil {
				return errOnResponse
			}
			return nil
		}

		payload := shared.RabbitMQPayload{
			Command: shared.AudioCommand,
			ID:      audio.ID,
		}

		audio.Status = shared.Preparing
		audio.ConvertTo = convertTo
		if errOnPrep := cc.cu.MarkAudioAsPreparing(audio); errOnPrep != nil {
			return errOnPrep
		}

		if errOnPublish := cc.RabbitMQ.PublishMessage(payload); errOnPublish != nil {
			return errOnPublish
		}

		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "⏳ Hold tight..., "+
			"your audio is being processed...")
		msg.ReplyToMessageID = int(audio.MessageID)
		if _, errOnSend := cc.TelegramAPI.Send(msg); errOnSend != nil {
			return errOnSend
		}
	}

	return nil
}
