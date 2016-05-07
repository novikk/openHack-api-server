package cron

import "github.com/robfig/cron"

func init() {
	c := cron.New()
	c.AddFunc("@every 10m", GetPollutionData)
	c.Start()
}
