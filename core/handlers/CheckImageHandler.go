package handlers

import (
	"IsUniquePhotoBot/core/repository"
	"IsUniquePhotoBot/core/util"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type CheckImageHandler struct {
	repo             repository.Repository
	remoteFileHasher *util.RemoteFileHasher
}

func NewCheckImageHandler(repo repository.Repository, remoteFileHasher *util.RemoteFileHasher) *CheckImageHandler {
	return &CheckImageHandler{
		repo:             repo,
		remoteFileHasher: remoteFileHasher,
	}
}

func (_ CheckImageHandler) Endpoint() any {
	return tele.OnPhoto
}

func (handler CheckImageHandler) Handle() tele.HandlerFunc {
	return handler.handleInternal
}

func (handler CheckImageHandler) handleInternal(ctx tele.Context) error {
	msg := ctx.Message()
	chatId := ctx.Chat().ID

	media := msg.Media()
	mediaType := media.MediaType()

	if mediaType != msg.Photo.MediaType() {
		ctx.Send("Cannot handle anything other than photos")
		return nil
	}

	fileId := media.MediaFile().FileID
	hash, err := handler.remoteFileHasher.DownloadAndHash(fileId)
	if err != nil {
		ctx.Reply("Failed to check photo")
	}

	groupId, err := handler.repo.GetGroupIdByUser(chatId)
	if err != nil {
		ctx.Reply("You didn't specify chat where to find photos")
	}

	existingMessageId, err := handler.repo.FindExisting(groupId, hash)
	if err != nil {
		ctx.Reply("Failed :c")
		return nil
	}

	if existingMessageId == -1 {
		ctx.Reply("It's unique!")
	} else {
		ctx.Reply(
			fmt.Sprintf(
				"Already [exists](%s)",
				util.MakeMessageLink(groupId, existingMessageId),
			),
		)

	}
	return nil
}

func sendFailureMessage(ctx tele.Context, message string) {
	ctx.Reply(message)

}
