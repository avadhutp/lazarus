package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/avadhutp/lazarus/ui"
	"github.com/codegangsta/cli"
	"github.com/gizak/termui"
	"github.com/go-ini/ini"
)

const (
	name    = "Lazarus"
	version = "0.0.1"
	desc    = "Lazarus: The resurrection of the revolution."
	hots    = "https://www.reddit.com/r/listentothis/hot.json?sort=hot"
)

// main Entry point for the application
func main() {
	log.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.Name = name
	app.Version = version
	app.Usage = desc
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "/etc/lazarus.ini",
			Usage: "An ini file with required configuration options. GO doesn't play well with home (~) paths. So specify the full path or use another absolute location.",
		},
	}

	app.Action = start
	app.Run(os.Args)
}

// start Starts all the necessary processes/go co-routines for the app to initialize
func start(ctx *cli.Context) {
	log.Info("Starting Lazarus...")
	cfg := checkAndReadConfig(ctx)

	ui.EventHandler()
	render(ctx)
	go downloadPlaylist(cfg)

	defer termui.Close()
	ui.Refresh()
	termui.Loop()
}

// checkAndReadConfig Loads config and makes sure it is ok
func checkAndReadConfig(ctx *cli.Context) *ui.Cfg {
	var cfg ui.Cfg
	if err := ini.MapTo(&cfg, ctx.String("config")); err != nil {
		panic(fmt.Sprintf("Problem loading the config file: %s", err.Error()))
	}

	if err := cfg.AllOk(); err != nil {
		panic(fmt.Sprintf("The ini file has problems: %s", err.Error()))
	}

	return &cfg
}

// render Paints the different widgest that compose Lazarus
func render(ctx *cli.Context) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	ui.Title.Label = fmt.Sprintf("*********** %s (%s) ***********", ctx.App.Name, ctx.App.Version)

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, ui.Title),
		),
		termui.NewRow(
			termui.NewCol(12, 0, ui.Songs),
		),
		termui.NewRow(
			termui.NewCol(6, 0, ui.Quit),
		),
	)
}

// downloadPlaylist Downloads the playlist from reddit and initiates the player
func downloadPlaylist(cfg *ui.Cfg) {
	player := ui.NewPlayer(cfg)

	player.Start(hots)
}
