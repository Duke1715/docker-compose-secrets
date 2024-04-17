package services

import (
	"docker-compose-secrets/app/client"
	"docker-compose-secrets/app/environment"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const GetSecretsDataPath = "v1/secrets/data"

type SecretService struct {
	httpClient         *client.HttpClient
	environmentService *environment.Service
}

func NewSecretService(
	httpClient *client.HttpClient, environmentService *environment.Service) *SecretService {
	return &SecretService{httpClient, environmentService}
}

func (s *SecretService) GetSecrets() (map[string]string, error) {
	headers := make(map[string]string)
	headers["X-Vault-Token"] = s.environmentService.GetVaultToken()

	response, err := s.httpClient.Get(
		s.buildUrl(),
		s.httpClient.BuildHeaders(headers),
	)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	if result["errors"] != nil {
		var errorMessages []string
		for _, e := range result["errors"].([]interface{}) {
			errorMessages = append(errorMessages, e.(string))
		}
		return nil, errors.New(strings.Join(errorMessages, "; "))
	}

	secrets := make(map[string]string)
	for secretKey, secretValue := range result {
		secrets[secretKey] = secretValue.(string)
	}

	return secrets, nil
}

func (s *SecretService) buildUrl() string {
	return fmt.Sprintf(
		"%s/%s/%s",
		s.environmentService.GetVaultAddr(),
		GetSecretsDataPath,
		s.environmentService.GetVaultPath(),
	)
}
