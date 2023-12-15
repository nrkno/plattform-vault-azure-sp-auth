package models

type VaultConfig struct {
	VaultAddress                   string
	VaultAzureRolesCredentialsPath string
}

type AzureCredentials struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
