package main

import (
	"flag"
	"fmt"
	"os"

	bot "github.com/joshbarrass/crossview-bot/internal"
	reddit "github.com/joshbarrass/goreddit"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// VERSION is the version of the bot
const VERSION = "0.0.1"

// USERAGENT is the user agent of the bot
const USERAGENT = "CrossView-Bot/%s"

// SUBREDDIT is the target subreddit
const SUBREDDIT = "crossview"

type Configuration struct {
	ClientID      string `envconfig:"REDDIT_CLIENT_ID" required:"true"`
	ClientSecret  string `envconfig:"REDDIT_CLIENT_SECRET" required:"true"`
	Username      string `envconfig:"REDDIT_USERNAME" required:"true"`
	Password      string `envconfig:"REDDIT_PASSWORD" required:"true"`
	DebugMode     bool   `envconfig:"DEBUG_MODE" default:"false"`
	ErrorUsername string `envconfig:"REDDIT_ERROR_REPORT_USERNAME"`
}

// handles flags before transferring to main
func init() {
	var getVersion bool
	flag.BoolVar(&getVersion, "version", false, "Get the current version")

	switch true {
	case getVersion:
		fmt.Printf("v%s\n", VERSION)
		os.Exit(0)
	}
}

func main() {
	var config Configuration
	err := envconfig.Process("", &config)
	if err != nil {
		logrus.Fatalf("could not process config: %s", err)
	}

	// format user agent with version
	userAgent := fmt.Sprintf(USERAGENT, VERSION)

	// create new reddit API
	redditAPI := reddit.API.NewRedditAPI(config.ClientID, config.ClientSecret, userAgent, config.Username, config.DebugMode)

	// authenticate
	err = redditAPI.Account.PasswordLogin(config.Password)
	if err != nil {
		logrus.Fatalf("unable to authenticate with reddit: %s", err)
	}
	logrus.WithFields(logrus.Fields{
		"token": fmt.Sprintf("%+v", redditAPI.Account.Token),
	}).Info("authenticated")

	// connect logging
	reddit.bot.SendErrors(redditAPI, config.ErrorUsername, userAgent)

	// initialise the bot
	// TODO: graceful shutdown?
	if err := bot.InitBot(SUBREDDIT); err != nil {
		logrus.Errorf("bot exited with error: %s", err)
	}
}
