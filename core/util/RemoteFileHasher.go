package util

import (
	"log"
	"os"

	tele "gopkg.in/telebot.v4"
)

type RemoteFileHasher struct {
	bot        *tele.Bot
	hashHelper *HashHelper
}

func NewRemoteFileHasher(bot *tele.Bot, hashHelper *HashHelper) *RemoteFileHasher {
	return &RemoteFileHasher{
		bot:        bot,
		hashHelper: hashHelper,
	}
}

func (rfh *RemoteFileHasher) DownloadAndHash(fileId string) (uint64, error) {
	photoPath := "./photos/photo.jpeg"

	enrichedFile, err := rfh.bot.FileByID(fileId)
	if err != nil {
		log.Printf("Error finding file %v", err)
	}
	err = rfh.bot.Download(&enrichedFile, photoPath)
	if err != nil {
		log.Printf("Error downloading file %v", err)
	}
	data, err := os.ReadFile(photoPath)
	if err != nil {
		log.Printf("Error readin %v", err)
		return 0, err
	}

	queryHash, err := rfh.hashHelper.computePHash(data)
	if err != nil {
		log.Printf("Error comptin hash: %v", err)
	}
	return queryHash, nil
}
