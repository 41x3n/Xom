package api

import (
	"github.com/41x3n/Xom/api/controller"
	"github.com/41x3n/Xom/core/domain"
	"github.com/41x3n/Xom/core/repository"
	"github.com/41x3n/Xom/core/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
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

func (h *rootHandler) HandleCallback(user *domain.User, callback *tgbotapi.CallbackQuery) error {
	telegramAPI := h.telegram.GetAPI()
	contextTimeout := time.Duration(h.env.ContextTimeout) * time.Second

	pr := repository.NewPhotoRepository(h.db, domain.TablePhoto)
	cu := usecase.NewCallbackUseCase(pr, contextTimeout)

	cc := controller.NewCallbackController(cu, telegramAPI)

	err := cc.HandleCallback(callback)

	return err
}
