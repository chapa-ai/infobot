package telegram

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"infoBot/internal/currencies"
	"infoBot/internal/db"
	"infoBot/internal/models"
	"log"
	"os"
	"sync"
)

var (
	bot   *tgbotapi.BotAPI
	bOnce sync.Once
)

func Run(logger *logrus.Entry) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := InitBot().GetUpdatesChan(u)
	err := UpdatesHandler(logger, updates)
	if err != nil {
		//log.Fatalln(err.Error())
		logrus.Errorf("failed handling updates: %s", err)
	}
}

func InitBot() *tgbotapi.BotAPI {
	bOnce.Do(func() {
		var err error
		bot, err = tgbotapi.NewBotAPI(os.Getenv("TelegramToken"))
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Println("Bot initialized")
	})
	return bot
}

func UpdatesHandler(logger *logrus.Entry, updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		message := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		err := handleMessages(logger, update, message)
		if err != nil {
			logrus.Errorf("failed handleMessages: %s", err)
			return err
		}
		_, err = InitBot().Send(message)
		if err != nil {
			logrus.Errorf("failed sending message: %s", err)
			return err
		}
	}
	return nil
}

func handleMessages(logger *logrus.Entry, update tgbotapi.Update, msg tgbotapi.MessageConfig) error {
	switch update.Message.Command() {
	case "start":
		message := tgbotapi.NewMessage(msg.ChatID, "Hello World")
		if _, err := bot.Send(message); err != nil {
			logger.Errorf("failed sending msg to telegram: %s", err)
			return err
		}
	case "info":
		message, err := GetInfo(logger, msg)
		if err != nil {
			logger.Errorf("failed GetInfo: %s", err)
			return err
		}
		if _, err := bot.Send(message); err != nil {
			logger.Errorf("failed sending msg to telegram: %s", err)
			return err
		}
	case "stat":
		message, err := GetStat(logger, msg)
		if err != nil {
			logger.Errorf("failed GetStat: %s", err)
			return err
		}
		if _, err := bot.Send(message); err != nil {
			logger.Errorf("failed sending msg to telegram: %s", err)
			return err
		}
	}
	return nil
}

func GetInfo(logger *logrus.Entry, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	currency, err := currencies.GetCurrency(logger, "BTC-USDT")
	if err != nil {
		logger.Errorf("failed currencies.GetCurrency: %s", err)
		return tgbotapi.MessageConfig{}, err
	}
	response := &models.Data{
		Symbol: currency.Data.Symbol,
		Buy:    currency.Data.Buy,
	}

	b, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("failed json.Marshal: %s", err)
		return tgbotapi.MessageConfig{}, err
	}

	_, err = db.SaveCurrencies(context.Background(), response)
	if err != nil {
		logger.Errorf("failed save currencies to db: %s", err)
		return tgbotapi.MessageConfig{}, err
	}

	return tgbotapi.NewMessage(msg.ChatID, string(b)), nil
}

func GetStat(logger *logrus.Entry, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	time, err := db.TimeOfFirstQuery(context.Background(), 1)
	if err != nil {
		logger.Errorf("failed get time of first query: %s", err)
		return tgbotapi.MessageConfig{}, err
	}

	count, err := db.CountOfAllQueries(context.Background())
	if err != nil {
		logger.Errorf("failed get count of all queries: %s", err)
		return tgbotapi.MessageConfig{}, err
	}

	response := &models.Stat{
		FirstQuery:        time,
		CountOfAllQueries: count,
	}

	b, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("failed json.Marshal: %s", err)
		return tgbotapi.MessageConfig{}, err
	}

	return tgbotapi.NewMessage(msg.ChatID, string(b)), nil
}
