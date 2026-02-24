package handlers

import tele "gopkg.in/telebot.v4"

type StartHandler struct{}

func (_ StartHandler) Endpoint() any {
	return "/start"
}

func (handler StartHandler) Handle() tele.HandlerFunc {
	return handler.handleInternal
}

func (_ StartHandler) handleInternal(ctx tele.Context) error {
	return ctx.Send("" +
		"Hi!\n" +
		"Forward me message from your channel\n" +
		"If your channel is private, then you have to add me as administrator\n" +
		"Then send me photos and I'll tell if any of them already in your channel\n",
	)
}
