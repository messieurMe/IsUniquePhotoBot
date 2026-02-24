package handlers

import (
	tele "gopkg.in/telebot.v4"
)

type Handler interface {
	Endpoint() any

	Handle() tele.HandlerFunc
}
