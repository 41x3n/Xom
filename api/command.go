package api

import (
	"time"

	"github.com/41x3n/Xom/api/controller"
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/core/repository"
	"github.com/41x3n/Xom/core/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *rootHandler) HandleStartCommand(user *domain.User, message *tgbotapi.Message) error {
	telegramAPI := h.telegram.GetAPI()

	sc := controller.NewStartController(telegramAPI)

	err := sc.HandleStartCommand(user, message)

	return err
}

func (h *rootHandler) HandleHelpCommand(user *domain.User, message *tgbotapi.Message) error {
	telegramAPI := h.telegram.GetAPI()

	hc := controller.NewHelpController(telegramAPI)

	err := hc.HandleHelpCommand(user, message)

	return err
}

func (h *rootHandler) HandlePhoto(user *domain.User, message *tgbotapi.Message) error {
	telegramAPI := h.telegram.GetAPI()
	contextTimeout := time.Duration(h.env.ContextTimeout) * time.Second

	pr := repository.NewPhotoRepository(h.db, domain.TablePhoto)
	pu := usecase.NewPhotoUsecase(pr, contextTimeout)
	pc := controller.NewPhotoController(pu, telegramAPI)

	err := pc.HandlePhoto(user, message)

	return err

}

func (h *rootHandler) HandleAudio(user *domain.User, message *tgbotapi.Message) error {
	telegramAPI := h.telegram.GetAPI()
	contextTimeout := time.Duration(h.env.ContextTimeout) * time.Second

	ar := repository.NewAudioRepository(h.db, domain.TableAudio)
	au := usecase.NewAudioUsecase(ar, contextTimeout)
	ac := controller.NewAudioController(au, telegramAPI)

	err := ac.HandleAudio(user, message)

	return err
}

func (h *rootHandler) HandleCallback(user *domain.User, callback *tgbotapi.CallbackQuery) error {
	telegramAPI := h.telegram.GetAPI()
	rabbitMQ := h.rabbitMQ
	contextTimeout := time.Duration(h.env.ContextTimeout) * time.Second

	pr := repository.NewPhotoRepository(h.db, domain.TablePhoto)
	ar := repository.NewAudioRepository(h.db, domain.TableAudio)
	cu := usecase.NewCallbackUseCase(pr, ar, contextTimeout)

	cc := controller.NewCallbackController(cu, telegramAPI, rabbitMQ)

	err := cc.HandleCallback(callback)

	return err
}
