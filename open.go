package main

import (
	aw "github.com/deanishe/awgo"
	"os"
)

func openChannel() {
	c := aw.NewCache(cache_dir)

	var c_channels []Channel
	if c.Exists(cache_file) {
		if err := c.LoadJSON(cache_file, &c_channels); err != nil {
			wf.FatalError(err)
		}

		for _, channel := range c_channels {
			if _, err := os.Stat(cache_dir + "/" + channel.ID + ".png"); err == nil {
				wf.NewItem(channel.Name).
					Var("teamID", channel.TeamID).
					Var("channelID", channel.ID).
					Icon(&aw.Icon{Value: cache_dir + "/" + channel.ID + ".png"}).
					Valid(true)
			} else {
				wf.NewItem(channel.Name).
					Var("teamID", channel.TeamID).
					Var("channelID", channel.ID).
					Valid(true)
			}
		}
	}

	args := wf.Args()
	if len(args) > 1 {
		wf.Filter(args[1])
	}

	wf.SendFeedback()
}
