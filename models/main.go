package models

type AzureConfig struct {
	ServicePrincipalPath string `json:"SERVICE_PRINCIPAL_PATH"`
	TenantID             string `json:"TENANT_ID"`
}

type AzureCredentials struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
