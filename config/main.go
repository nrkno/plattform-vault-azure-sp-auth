package config

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/nrkno/plattform-vault-azure-sp-auth/utils"

	vault "github.com/hashicorp/vault/api"
)

var DefaultReadVaultPathOptions = ReadVaultPathOptions{
	RetryCount: utils.ToPointer(5),
	RetryTime:  utils.ToPointer(3 * time.Second),
	Logger:     nil,
}

type ReadVaultPathOptions struct {
	// The amount of times to retry before returning (default 5)
	RetryCount *int
	// The duration to wait before next retry (default 3 * time.Second)
	RetryTime *time.Duration
	// Logger for calling Warn(err) fails. Last fail will be returned as error
	Logger *log.Logger
}

func ReadVaultPath[config any](vaultAddress string, path string, opts *ReadVaultPathOptions) (*config, error) {
	if opts == nil {
		opts = &DefaultReadVaultPathOptions
	} else {
		if opts.RetryCount == nil {
			opts.RetryCount = DefaultReadVaultPathOptions.RetryCount
		}
		if opts.RetryTime == nil {
			opts.RetryTime = DefaultReadVaultPathOptions.RetryTime
		}
	}

	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	vaultClient, err := vault.NewClient(&vault.Config{Address: vaultAddress, HttpClient: httpClient})
	if err != nil {
		return nil, err
	}

	var secret *vault.Secret
	for i := 0; i <= *opts.RetryCount; i++ {
		secret, err = vaultClient.Logical().Read(path)
		if err != nil {
			opts.Logger.Fatal(err.Error())
			time.Sleep(*opts.RetryTime)
		} else {
			break
		}
	}

	if secret == nil {
		secret, err = vaultClient.Logical().Read(path)
		if err != nil {
			return nil, err
		}
	}

	data, err := json.Marshal(secret.Data)
	if err != nil {
		return nil, err
	}

	var conf config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
