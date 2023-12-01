package GetCreds

import (
	"log/syslog"
	"os"

	"github.com/nrkno/plattform-vault-azure-sp-auth/config"
	"github.com/nrkno/plattform-vault-azure-sp-auth/models"
)

func GetCreds(VAULT_ADDR string, VAULT_AZURE_CONFIG_PATH string) (*models.AzureCredentials, error) {
	log, err := syslog.New(syslog.LOG_INFO|syslog.LOG_USER, os.Args[0])
	if err != nil {
		return nil, err
	}
	defer log.Close()

	// Vault Azure SP
	var conf *models.AzureConfig
	conf, err = config.ReadVaultPath[models.AzureConfig](VAULT_ADDR, VAULT_AZURE_CONFIG_PATH,
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
