package main

import (
	"healthCheck/internal/app"
	"log"
	"time"
)

func main() {
	log.Println("starting application...")

	// todo: os.Interrupt
	config, err := app.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(config)
	for {
		if err := app.Run(); err != nil {
			app.ReportError(err)
		}

		app.Log.Infoln("ok")
		time.Sleep(time.Duration(config.RetryAfter) * time.Second)
	}
}
