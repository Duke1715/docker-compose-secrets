package models

type SecretData struct {
	Data map[string]string `json:"data"`
}

type GetSecretsResult struct {
	Errors []string    `json:"errors,omitempty"`
	Data   *SecretData `json:"data,omitempty"`
}
