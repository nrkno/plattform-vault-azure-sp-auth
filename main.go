package GetCreds

import (
	"log"

	"github.com/nrkno/plattform-vault-azure-sp-auth/config"
	"github.com/nrkno/plattform-vault-azure-sp-auth/models"
)

func GetCreds(VAULT_ADDR string, VAULT_AZURE_CONFIG_PATH string) (*models.AzureCredentials, error) {
	log := &log.Logger{}

	// Vault Azure SP
	var conf *models.AzureConfig
	conf, err := config.ReadVaultPath[models.AzureConfig](VAULT_ADDR, VAULT_AZURE_CONFIG_PATH,
		&config.ReadVaultPathOptions{Logger: log})
	if err != nil {
		return nil, err
	}
	spCred, err := config.ReadVaultPath[models.AzureCredentials](VAULT_ADDR, conf.ServicePrincipalPath,
		&config.ReadVaultPathOptions{Logger: log})
	if err != nil {
		return nil, err
	}

	return spCred, nil
}
