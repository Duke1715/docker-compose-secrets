package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Get the arguments passed to the program
	args := os.Args[1:]

	// Check if the user passed any arguments
	if len(args) == 0 {
		invalidCommand()
	}

	// Check if the required environment variables are set
	if os.Getenv("VAULT_ADDR") == "" || os.Getenv("VAULT_TOKEN") == "" {
		log.Fatal("VAULT_ADDR and VAULT_TOKEN must be set in environment")
	}

	// Check which command the user passed
	switch args[0] {
	case "start":
		start(false)
	case "stop":
		stop()
	case "restart":
		start(true)
	case "update":
		update()
	default:
		invalidCommand()
	}
}

func invalidCommand() {
	fmt.Println("Invalid command, please use one of the following:")
	fmt.Println("  start")
	fmt.Println("  stop")
	fmt.Println("  restart")
	fmt.Println("  update")
	fmt.Println("\nExample: dcs start")
	os.Exit(1)
}

func start(restart bool) {
	// Get the address of the Vault server from the environment
	server := os.Getenv("VAULT_ADDR")

	// Get the folder of the current working directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Base(dir)

	// Get the token from the environment
	token := os.Getenv("VAULT_TOKEN")

	// Print the settings to the user
	fmt.Println("Retrieving secrets from Vault:")
	fmt.Println("  Server:", server)
	fmt.Println("  Path:", path)
	fmt.Println("  Token:", token)
	fmt.Println()

	// Initialize the request to the Vault server with the correct path
	req, _ := http.NewRequest("GET", server+"/v1/secret/data/"+path, nil)

	// Add the authentication token header to the request
	req.Header.Add("X-Vault-Token", token)

	// Send the request to the Vault server
	res, _ := http.DefaultClient.Do(req)

	// Close the response body when we're done
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	// Decode the response body
	body, _ := io.ReadAll(res.Body)

	// Unmarshal and parse the JSON response into a map
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	// Set the command to run depending on if we're restarting or not
	var cmd *exec.Cmd
	if restart {
		cmd = exec.Command("docker", "compose", "up", "-d", "--force-recreate")
	} else {
		cmd = exec.Command("docker", "compose", "up", "-d")
	}

	// secrets := result["data"].(map[string]interface{})["data"].(map[string]interface{})

	// Extract the secrets from the response and handle errors
	s1 := result["data"]
	if s1 != nil {
		s2 := s1.(map[string]interface{})
		if s2 != nil {
			s3 := s2["data"]
			if s3 != nil {
				s4 := s3.(map[string]interface{})
				if s4 != nil {
					fmt.Println("Injecting secrets into process:")

					// Pass all OS environment variables to the command
					cmd.Env = os.Environ()

					// Inject all secrets into the command as environment variables and print them to the user
					for k, v := range s4 {
						cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=%s", k, v))
						fmt.Printf("  %s: %s\n", k, v)
					}
				} else {
					fmt.Println("No secrets found for \"" + path + "\", continuing without secrets")
				}
			} else {
				fmt.Println("No secrets found for \"" + path + "\", continuing without secrets")
			}
		} else {
			fmt.Println("No secrets found for \"" + path + "\", continuing without secrets")
		}
	} else {
		fmt.Println("No secrets found for \"" + path + "\", continuing without secrets")
	}

	fmt.Println()
	fmt.Println("Starting docker compose:")

	// Write the output of the command to the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func stop() {
	// Set the command to run
	cmd := exec.Command("docker", "compose", "down", "--remove-orphans")

	// Pass all OS environment variables to the command
	cmd.Env = os.Environ()

	fmt.Println("Stopping docker compose:")

	// Write the output of the command to the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func update() {
	// Set the command to run
	cmd := exec.Command("docker", "compose", "pull")

	// Pass all OS environment variables to the command
	cmd.Env = os.Environ()

	fmt.Println("Pulling latest docker images:")

	// Write the output of the command to the terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	// Restart the docker compose after updating the images
	start(true)
}
