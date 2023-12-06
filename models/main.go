package models

type VaultConfig struct {
	VaultAddress         string
	VaultAzureConfigPath string
}

type AzureCredentials struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
