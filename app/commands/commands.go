package commands

import (
	"docker-compose-secrets/app/services"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	CommandStart   = "start"
	CommandStop    = "stop"
	CommandRestart = "restart"
	CommandUpdate  = "update"
)

type Service struct {
	secretService     *services.SecretService
	availableCommands []string
	currentCommand    string
}

func NewService(secretService *services.SecretService) *Service {
	return &Service{
		secretService: secretService,
		availableCommands: []string{
			CommandStart,
			CommandStop,
			CommandRestart,
			CommandUpdate,
		},
	}
}

func (c *Service) CheckCommandName(argCommand string) error {
	for _, cmd := range c.availableCommands {
		if cmd == argCommand {
			return nil
		}
	}

	return fmt.Errorf("invalid command `%s`", argCommand)
}

func (c *Service) SetCurrentCommand(argCommand string) {
	c.currentCommand = argCommand
}

func (c *Service) GetCurrentCommand() string {
	return c.currentCommand
}

func (c *Service) ExecuteStart() {
	secrets, err := c.secretService.GetSecrets()
	if err != nil {
		log.Fatal(err)
	}

	cmd := buildCommand("compose", "up", "-d")
	addEnvironToCmd(cmd, secrets)

	println("Starting docker compose:")
	if err := executeCommand(cmd); err != nil {
		log.Fatal(err)
	}
}

func (c *Service) ExecuteStop() {
	println("stopping docker compose:")

	if err := executeCommand(
		buildCommand("compose", "down", "--remove-orphans"),
	); err != nil {
		log.Fatal(err)
	}
}

func (c *Service) ExecuteRestart() {
	secrets, err := c.secretService.GetSecrets()
	if err != nil {
		log.Fatal(err)
	}

	cmd := buildCommand("compose", "up", "-d", "--force-recreate")
	addEnvironToCmd(cmd, secrets)

	println("Starting docker compose:")
	if err := executeCommand(cmd); err != nil {
		log.Fatal(err)
	}
}

func (c *Service) ExecuteUpdate() {
	println("pulling latest docker images:")

	if err := executeCommand(
		buildCommand("compose", "pull"),
	); err != nil {
		log.Fatal(err)
	}
}

func buildCommand(args ...string) *exec.Cmd {
	cmd := exec.Command("docker", args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func addEnvironToCmd(cmd *exec.Cmd, env map[string]string) {
	if len(env) == 0 {
		println("No secrets found. Continuing without secrets")
		return
	}

	for key, value := range env {
		cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=%s", key, value))
	}
}

func executeCommand(cmd *exec.Cmd) error {
	return cmd.Run()
}
