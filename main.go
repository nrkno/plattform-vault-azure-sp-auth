package GetCreds

import (
	"github.com/nrkno/plattform-vault-azure-sp-auth/config"
	"github.com/nrkno/plattform-vault-azure-sp-auth/models"
)

// VaultReader defines the interface for reading from Vault
type VaultReader interface {
	ReadVaultPath(vaultCfg models.VaultConfig, opts *config.ReadVaultPathOptions) (*models.AzureCredentials, error)
}

// realVaultReader wraps the actual config.ReadVaultPath implementation
type realVaultReader struct{}

func (r *realVaultReader) ReadVaultPath(vaultCfg models.VaultConfig, opts *config.ReadVaultPathOptions) (*models.AzureCredentials, error) {
	return config.ReadVaultPath[models.AzureCredentials](vaultCfg, opts)
}

// GetCreds reads Azure credentials from Vault using the provided reader
func GetCreds(reader VaultReader, VAULT_ADDR string, VAULT_AZURE_ROLES_CREDENTIALS_PATH string) (*models.AzureCredentials, error) {

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
	spCred, err := reader.ReadVaultPath(vaultCfg, nil)
	if err != nil {
		return nil, err
	}

	return spCred, nil
}

// GetCredsWithDefaults is a convenience function that uses the real Vault reader
func GetCredsWithDefaults(VAULT_ADDR string, VAULT_AZURE_ROLES_CREDENTIALS_PATH string) (*models.AzureCredentials, error) {
	return GetCreds(&realVaultReader{}, VAULT_ADDR, VAULT_AZURE_ROLES_CREDENTIALS_PATH)
}
