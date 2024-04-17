package environment

import (
	"errors"
	"fmt"
	"os"
)

const (
	envVaultAddr  = "VAULT_ADDR"
	envVaultToken = "VAULT_TOKEN"
	envVaultPath  = "VAULT_PATH"
)

type Service struct {
	systemEnv map[string]string
}

func NewService() *Service {
	systemEnv := make(map[string]string)
	systemEnv[envVaultAddr] = os.Getenv(envVaultAddr)
	systemEnv[envVaultToken] = os.Getenv(envVaultToken)
	systemEnv[envVaultPath] = os.Getenv(envVaultPath)

	return &Service{systemEnv}
}

func (e *Service) CheckExistSystemEnv() error {
	for envKey, envVal := range e.systemEnv {
		if envVal == "" {
			return errors.New(
				fmt.Sprintf(
					"environment variable `%s` is empty",
					envKey,
				),
			)
		}
	}

	return nil
}

func (e *Service) GetVaultAddr() string {
	return e.systemEnv[envVaultAddr]
}

func (e *Service) GetVaultToken() string {
	return e.systemEnv[envVaultToken]
}

func (e *Service) GetVaultPath() string {
	return e.systemEnv[envVaultPath]
}
