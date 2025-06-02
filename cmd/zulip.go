package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/wakumaku/go-zulip"
	"github.com/wakumaku/go-zulip/realtime"
	"github.com/wakumaku/go-zulip/realtime/events"
)

// Check user in blacklist
func inBlackList(user User, senderEmail string) bool {
	for _, from := range user.IgnorePrivateFrom {
		if from == senderEmail {
			return true
		}
	}

	return false
}

// Find word form words list in message
func hasWord(user User, message string) bool {
	for _, word := range user.Words {
		if strings.Contains(message, word) {
			return true
		}
	}

	return false
}

// Send message to Telegram bot
func sendMessageToTG(user User, messageTG string) error {
	log.Debugf("message Zulip->TG: %s", messageTG)

	msg := tgbotapi.NewMessage(user.ID, messageTG)
	msg.LinkPreviewOptions.IsDisabled = true
	msg.ParseMode = "Markdown"

	_, err := botTG.Send(msg)
	if err != nil {
		log.Errorf("Telegram message send failed: %v", err)

		// attempt to send again without Markdown
		msg.ParseMode = ""
		_, err = botTG.Send(msg)
		if err != nil {
			log.Errorf("Telegram message send failed: %v", err)
			return err
		}
	}

	return nil
}

func taskZulip(ctx context.Context, user User) {
	attempt := 0

	for {
		if attempt > 1 {
			time.Sleep(time.Second * 30)
		}
		attempt++

		log.Debugf("Connection attempt %d for user %s", attempt, user.Name)

		// Создаем новый экземпляр клиента Zulip
		credentials := zulip.Credentials(config.Zulip.Site, user.ZulipEmail, user.ZulipKey)
		clientZulip, err := zulip.NewClient(credentials)
		if err != nil {
			log.Errorf("Error creating Zulip client for user %s: %v", user.Name, err)
			continue
		}
		log.Infof("Zulip client %s successfully launched with attempt %d", user.Name, attempt)

		realtimeSvc := realtime.NewService(clientZulip)

		queue, err := realtimeSvc.RegisterEvetQueue(ctx,
			realtime.EventTypes(
				events.MessageType,
			),
			realtime.AllPublicStreams(true),
		)
		if err != nil {
			log.Errorf("Error registering event queue for user %s: %s", user.Name, err)
			continue
		}

		if queue.IsError() {
			log.Errorf("%s: %s for user %s", queue.Msg(), queue.Code(), user.Name)
			continue
		}

		log.Debugf("QueueId: %s for user %s", queue.QueueId, user.Name)

		lastEventID := queue.LastEventId

		// Infinite loop polling for new events
		for {
			// Long polling HTTP Request
			eventsFromQueue, err := realtimeSvc.GetEventsEventQueue(ctx, queue.QueueId, realtime.LastEventID(lastEventID))
			if err != nil {
				log.Errorf("error getting events from queue: %s", err)
				break
			}

			for _, e := range eventsFromQueue.Events {
				lastEventID = e.EventID()

				if e.EventType() == events.MessageType {
					message := e.(*events.Message)
					messageHeader := ""
					messageContent := message.Message.Content

					if message.Message.DisplayRecipient.IsChannel {
						channel := message.Message.DisplayRecipient.Channel
						topic := message.Message.Subject
						messageHeader = fmt.Sprintf("*%s:%s %s*\n", channel, topic, message.Message.SenderFullName)

						if !hasWord(user, messageContent) {
							continue
						}
					} else { // private message
						if inBlackList(user, message.Message.SenderEmail) {
							continue
						}

						messageHeader = fmt.Sprintf("*%s*\n", message.Message.SenderFullName)
					}

					messageTG := messageHeader + messageContent
					sendMessageToTG(user, messageTG)
				}
			}
		}

		log.Errorf("Break main loop for user %s", user.Name)
	}
}
