package handlers

import (
	"IsUniquePhotoBot/core/repository"
	"IsUniquePhotoBot/core/util"
	"log"
	"strconv"

	tele "gopkg.in/telebot.v4"
)

type ScanGroupHandler struct {
	repo             repository.Repository
	bot              *tele.Bot
	remoteFileHasher *util.RemoteFileHasher
}

func NewScanGroupHandler(
	repo repository.Repository,
	bot *tele.Bot,
	remoteFileHasher *util.RemoteFileHasher,
) *ScanGroupHandler {
	return &ScanGroupHandler{
		repo:             repo,
		bot:              bot,
		remoteFileHasher: remoteFileHasher,
	}
}

func (_ ScanGroupHandler) Endpoint() any {
	return tele.OnForward
}

func (handler ScanGroupHandler) Handle() tele.HandlerFunc {
	return handler.handleInternal
}

func (handler ScanGroupHandler) handleInternal(ctx tele.Context) error {
	msg := ctx.Message()
	currentChatId := ctx.Chat().ID
	originChatId := msg.Origin.Chat.ID

	currentMax, _ := handler.repo.GetLastGroupMessageByUser(currentChatId)

	maxId := handler.scanGroup(ctx, originChatId, currentChatId, currentMax)

	handler.repo.SetUserGroupAndLastMessage(
		currentChatId,
		originChatId,
		maxId,
	)
	return nil
}

func (handler ScanGroupHandler) scanGroup(
	ctx tele.Context,
	originChatId int64,
	chatId int64,
	existingMaxId int,
) int {
	id := 0
	notFoundCounter := 0
	tryForwardWindow := 15

	if existingMaxId > id {
		id = existingMaxId
	}

	for notFoundCounter < tryForwardWindow {
		id++

		foundImage := handler.forwardProbeAndDelete(ctx, chatId, originChatId, id)

		if foundImage {
			notFoundCounter = 0
		} else {
			notFoundCounter++
		}
	}

	if id == tryForwardWindow {
		ctx.Send("Sorry, I failed to get images from your chat")
		return 0
	}
	return id - tryForwardWindow
}

func (handler ScanGroupHandler) forwardProbeAndDelete(
	ctx tele.Context,
	chatId int64,
	originChatId int64,
	messageId int,
) bool {
	msg, err := handler.bot.Forward(ctx.Recipient(), newEditableImpl(messageId, originChatId))

	if err != nil {
		log.Printf("Failed to forward %v", err)
		return false
	}

	defer handler.bot.Delete(msg)

	hash := handler.probeMessage(msg)
	if hash == nil {
		return false
	}

	handler.repo.AddImage(
		chatId,
		msg.Origin.MessageID,
		*hash,
	)
	return true
}

func (handler ScanGroupHandler) probeMessage(msg *tele.Message) *uint64 {
	if msg == nil {
		return nil
	}

	var fileID string
	if msg.Photo != nil {
		fileID = msg.Photo.FileID
	} else {
		return nil
	}

	hash, err := handler.remoteFileHasher.DownloadAndHash(fileID)

	log.Printf("HASH: %s : %v", fileID, hash)

	if err != nil {
		return nil
	}

	return &hash
}

type EditableImpl struct {
	messageId string
	chatId    int64
}

func newEditableImpl(messageId int, chatId int64) EditableImpl {
	return EditableImpl{
		messageId: strconv.Itoa(messageId),
		chatId:    chatId,
	}
}

func (editable EditableImpl) MessageSig() (messageID string, chatID int64) {
	return editable.messageId, editable.chatId
}
