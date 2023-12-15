package GetCreds

import (
	"github.com/nrkno/plattform-vault-azure-sp-auth/config"
	"github.com/nrkno/plattform-vault-azure-sp-auth/models"
)

func GetCreds(VAULT_ADDR string, VAULT_AZURE_ROLES_CREDENTIALS_PATH string) (*models.AzureCredentials, error) {

	if VAULT_ADDR == "" {
		VAULT_ADDR = "http://localhost:8200"
	}
	if VAULT_AZURE_ROLES_CREDENTIALS_PATH == "" {
		VAULT_AZURE_ROLES_CREDENTIALS_PATH = "default_vault_azure_roles_path"
	}
	// Vault Azure SP
	var vaultCfg models.VaultConfig = models.VaultConfig{
		VaultAddress:                   VAULT_ADDR,
		VaultAzureRolesCredentialsPath: VAULT_AZURE_ROLES_CREDENTIALS_PATH,
	}
	spCred, err := config.ReadVaultPath[models.AzureCredentials](vaultCfg, nil)
	if err != nil {
		return nil, err
	}

	return spCred, nil
}
