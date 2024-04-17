package services

import (
	"docker-compose-secrets/app/client"
	"docker-compose-secrets/app/environment"
	"docker-compose-secrets/app/models"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const getSecretsDataPath = "v1/secret/data"

type SecretService struct {
	httpClient         *client.HttpClient
	environmentService *environment.Service
}

func NewSecretService(
	httpClient *client.HttpClient, environmentService *environment.Service,
) *SecretService {
	return &SecretService{httpClient, environmentService}
}

func (s *SecretService) GetSecrets() (map[string]string, error) {
	headers := make(map[string]string)
	headers[client.HeaderVaultTokenName] = s.environmentService.GetVaultToken()

	response, err := s.httpClient.Get(
		s.buildUrl(),
		s.httpClient.BuildHeaders(headers),
	)
	if err != nil {
		return nil, err
	}

	var result models.GetSecretsResult
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, errors.New(strings.Join(result.Errors, "; "))
	}

	secrets := make(map[string]string)
	for secretKey, secretValue := range result.Data.Data {
		secrets[secretKey] = secretValue
	}

	return secrets, nil
}

func (s *SecretService) buildUrl() string {
	return fmt.Sprintf(
		"%s/%s/%s",
		s.environmentService.GetVaultAddr(),
		getSecretsDataPath,
		s.environmentService.GetVaultPath(),
	)
}
