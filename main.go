package main

import (
	"database/sql"
	"log"
	"os"
	"scroogebot/expenditure"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

const (
	unknownCommandText       string = "I don't know that command"
	messageHandlingErrorText string = "Sorry! The last message did not handled :("
	calculationErrorText     string = "Sorry! The calculation not working :("
)

func main() {
	db := prepareDb()
	defer db.Close()

	logfile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer logfile.Close()
	infoLog := log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	token := os.Getenv("SCROOGE_TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "calc":
				excuteCalcCommand(bot, db, errorLogger, update.Message.Chat.ID, time.Now().AddDate(0, 0, -time.Now().Day()))
			case "calc_week":
				excuteCalcCommand(bot, db, errorLogger, update.Message.Chat.ID, time.Now().AddDate(0, 0, -7))
			case "calc_month":
				excuteCalcCommand(bot, db, errorLogger, update.Message.Chat.ID, time.Now().AddDate(0, 0, -30))
			case "calc_quarter":
				excuteCalcCommand(bot, db, errorLogger, update.Message.Chat.ID, time.Now().AddDate(0, 0, -90))
			default:
				sendErrorMessage(bot, errorLogger, update.Message.Chat.ID, unknownCommandText)
			}
		} else if newRecord, err := expenditure.Parse(update.Message.Text); newRecord != nil {

			if err != nil {
				infoLog.Println(err)
				sendErrorMessage(bot, errorLogger, update.Message.Chat.ID, err.Error())
				continue
			}

			newRecord.SetContextData(update.Message.Chat.ID, update.Message.From.UserName)

			if dbError := newRecord.Save(db); dbError != nil {
				errorLogger.Println(err.Error())
				sendErrorMessage(bot, errorLogger, update.Message.Chat.ID, messageHandlingErrorText)
			}
		}
	}
}

func excuteCalcCommand(bot *tgbotapi.BotAPI, db *sql.DB, logger *log.Logger, chatId int64, from time.Time) {
	msg := tgbotapi.NewMessage(chatId, "")

	if results, err := expenditure.GetRecords(db, chatId, from); err != nil {
		sendErrorMessage(bot, logger, chatId, calculationErrorText)
	} else {
		msg.Text = expenditure.MakeCalculationResult(results)
		if _, err := bot.Send(msg); err != nil {
			logger.Println(err)
		}
	}
}

func sendErrorMessage(bot *tgbotapi.BotAPI, logger *log.Logger, chatId int64, errorText string) {
	msg := tgbotapi.NewMessage(chatId, errorText)

	if _, err := bot.Send(msg); err != nil {
		logger.Println(err.Error())
	}
}

func prepareDb() *sql.DB {
	_, error := os.Stat("scrooge.db")
	fileNotExist := os.IsNotExist(error)
	if fileNotExist {
		file, err := os.Create("scrooge.db")
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", "scrooge.db")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if fileNotExist {
		expenses_table := `CREATE TABLE "Expenses" ( 
			"Id" INTEGER NOT NULL, 
			"Category" TEXT NOT NULL, 
			"Amount" REAL NOT NULL, 
			"ChatId" INTEGER NOT NULL, 
			"CreatedAt" TEXT NOT NULL, 
			"CreatedBy" TEXT NOT NULL, 
			PRIMARY KEY("Id" AUTOINCREMENT) )`
		query, err := db.Prepare(expenses_table)
		if err != nil {
			log.Fatal(err)
		}
		query.Exec()
	}

	return db
}
