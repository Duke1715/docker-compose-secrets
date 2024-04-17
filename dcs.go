package main

import (
	"docker-compose-secrets/app"
	"docker-compose-secrets/app/client"
	"docker-compose-secrets/app/commands"
	"docker-compose-secrets/app/environment"
	"docker-compose-secrets/app/services"
	"log"
	"os"
)

func main() {
	environmentService := environment.NewService()
	if err := environmentService.CheckExistSystemEnv(); err != nil {
		log.Fatal(err)
	}

	args := os.Args[1:]

	commandService := commands.NewService(
		services.NewSecretService(
			client.NewHttpClient(),
			environmentService,
		),
	)
	checkArgument(args, commandService)
	commandService.SetCurrentCommand(args[0])

	application := app.NewApplication(commandService)
	application.Run()
}

func checkArgument(args []string, commands *commands.Service) {
	if len(args) == 0 {
		showWrongCommandError("invalid command")
	}

	if err := commands.CheckCommandName(args[0]); err != nil {
		showWrongCommandError(err.Error())
	}
}

func showWrongCommandError(title string) {
	log.Fatalf(
		"%s, please use one of the following:\n %s\n %s\n %s\n %s\n%s",
		title,
		"start",
		"stop",
		"restart",
		"update",
		"Example: dcs start",
	)
}
