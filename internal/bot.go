package bot

import reddit "github.com/joshbarrass/goreddit/API"

type Bot struct {
	API       *reddit.RedditAPI
	Subreddit string
	DebugMode bool
}

func InitBot(api *reddit.RedditAPI, subreddit string, debugMode bool) error {
	bot := Bot{api, subreddit, debugMode}
	return bot.Main()
}

// Main is the main method of the bot
func (bot *Bot) Main() error {

}
