package main

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	_ "Nosviak4/source/functions/iplookup"
	"Nosviak4/source/masters"
	"Nosviak4/source/masters/commands"
	_ "Nosviak4/source/masters/commands/commands"
	_ "Nosviak4/source/masters/commands/subcommands"
	"Nosviak4/source/web"

	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
)

// inits all the components
func main() {
	if err := source.LOGGER.SetAggregatedLogger(filepath.Join(source.ASSETS, "logs", "terminal.txt"), 100000); err != nil {
		log.Panic(err)
	}

	term := source.LOGGER.AggregateTerminal()
	gologr.DEBUGENABLED = slices.Contains(os.Args, "-d")
	term.WriteLog(gologr.DEFAULT, "Initializing Nosviak4 %s", source.VERSION)
	if err := source.OpenOptions(); err != nil {
		term.WriteLog(gologr.ERROR, "Error occurred while parsing Config: %v", err)
		return
	}

	/* defines the params for the root args */
	commands.ROOT.Args = append(make([]*commands.Arg, 0), commands.Descriptor, commands.Target, commands.Duration, commands.Port)
	if source.OPTIONS.Bool("attacks", "port_then_duration") {
		commands.ROOT.Args = append(make([]*commands.Arg, 0), commands.Descriptor, commands.Target, commands.Port, commands.Duration)
	}

	term.WriteLog(gologr.DEFAULT, "[Successfully initialized the GoConfig interface] (%d)", len(source.OPTIONS.Config.Renders))
	if err := database.DB.Connect(); err != nil {
		term.WriteLog(gologr.ERROR, "Error occurred while connecting to the database: %v", err)
		return
	}

	term.WriteLog(gologr.DEFAULT, "[Successfully connected to the database] (%s)", database.DB.Config.Database)
	if err := web.NewWebServe(); err != nil {
		term.WriteLog(gologr.ERROR, "Error occurred while starting web server: %v", err)
		return
	}

	if err := commands.OpenBotConn(); err != nil {
		term.WriteLog(gologr.ERROR, "Error occurred while starting the bot: %v", err)
		return
	}

	master, err := masters.NewMaster()
	if err != nil {
		term.WriteLog(gologr.ERROR, "Error occurred while creating the server: %v", err)
		return
	}

	if err := master.Listen(); err != nil {
		term.WriteLog(gologr.ERROR, "Error occurred while listening for conns: %v", err)
		return
	}
}
