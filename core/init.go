package core

import (
	"IsUniquePhotoBot/core/handlers"
	"IsUniquePhotoBot/core/repository/sqlite"
	"IsUniquePhotoBot/core/util"
	"database/sql"
	"log"

	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	tele "gopkg.in/telebot.v4"
)

func StartBot() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("BOT_KEY")

	pref := tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	db, err := sql.Open("sqlite3", "db/mybot.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	similarityThreshold := 5
	hashHelper := util.NewHashHelper(similarityThreshold)

	repository, err := sqlite.NewSQLiteRepository(db, hashHelper)

	remoteFileHasher := util.NewRemoteFileHasher(bot, hashHelper)

	handlers := []handlers.Handler{
		handlers.StartHandler{},
		handlers.NewScanGroupHandler(repository, bot, remoteFileHasher),
		handlers.NewCheckImageHandler(repository, remoteFileHasher),
	}

	for _, handler := range handlers {
		bot.Handle(
			handler.Endpoint(),
			handler.Handle(),
		)
	}
	bot.Start()
}
