package app

import (
	"docker-compose-secrets/app/commands"
	"log"
)

type Application struct {
	commandService *commands.Service
}

func NewApplication(commandService *commands.Service) *Application {
	return &Application{commandService}
}

func (app *Application) Run() {
	currentCommand := app.commandService.GetCurrentCommand()
	if currentCommand == "" {
		log.Fatal("no command specified")
	}

	switch currentCommand {
	case commands.CommandStart:
		app.commandService.ExecuteStart()
	case commands.CommandStop:
		app.commandService.ExecuteStop()
	case commands.CommandRestart:
		app.commandService.ExecuteRestart()
	case commands.CommandUpdate:
		app.commandService.ExecuteUpdate()
		app.commandService.ExecuteRestart()
	}
}
