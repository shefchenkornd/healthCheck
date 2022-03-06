package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

// App is base instance
type App struct {
	Config *Config
	Log    logrus.FieldLogger
}

// NewApp is App constructor
func NewApp(config *Config) *App {
	log := logrus.New()

	// Logrus timestamp formatting
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	logFile, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("can't create log file, err: %v", err)
	}
	log.Out = logFile

	return &App{
		Config: config,
		Log:    log,
	}
}

// Run perform a health check on a specific URL
func (app App) Run() error {
	resp, err := http.Head(app.Config.URL)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v site returned HTTP status %v", app.Config.URL, resp.Status)
	}

	return nil
}

// ReportError report about error
func (app App) ReportError(err error) {
	msgErr := fmt.Errorf(
		"URL: %v \nError: %v\nTime: %v",
		app.Config.URL,
		err,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	app.Log.Error(msgErr)

	// send message to telegram
	if errTg := app.sendToTelegram(msgErr.Error()); errTg != nil {
		app.Log.Error(errTg)
	}
}

// sendToTelegram send message to telegram chat
func (app App) sendToTelegram(msg string) error {
	msg = strings.ReplaceAll(msg, "\n", "%0A")
	msg = strings.ReplaceAll(msg, "https://", "")
	endpoint := fmt.Sprintf(
		"https://api.telegram.org/bot%v/sendMessage?chat_id=%v&text=%v",
		app.Config.Telegram.Token,
		app.Config.Telegram.ChatId,
		msg,
	)
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("can't send message to telegram")
	}

	return nil
}
