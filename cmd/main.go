package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

const FILE_LOG = "bot.log"

var (
	config Config
	botTG  *tgbotapi.BotAPI
)

func main() {
	var err error
	ctx := context.Background()

	// Open log file
	fileLog, _ := os.OpenFile(FILE_LOG, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer fileLog.Close()

	log.SetOutput(fileLog)
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% [%lvl%] %msg%\n",
	})

	log.SetLevel(log.DebugLevel)
	log.Info("Start")

	// Read config
	err = readConfig()
	if err != nil {
		log.Fatalf("Error config: %v", err)
		os.Exit(1)
	}

	// Create Telegram bot
	botTG, err = tgbotapi.NewBotAPI(config.TGKey)
	if err != nil {
		log.Fatalf("Error Telegram bot launch: %v", err)
		os.Exit(1)
	}

	botTG.Debug = false // Debug mode (optional)
	log.Infof("Telegram bot %s successfully launched!", botTG.Self.UserName)

	// Channel for getting signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for _, user := range config.Users {
		go taskZulip(ctx, user)
	}

	<-c

	log.Info("Finish")
}
