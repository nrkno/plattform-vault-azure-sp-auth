# Go package for dynamically generating temporary Azure Service Principals 
This package integrates with Vault to automatically generate temporary Service Principals in Azure, which may used as authentication by your Go code.

## Prerequisites
1. You need to define the level of access for your service principal in the Vault definitions
1. The Vault definitions needs to be approved and applied
1. 

## How to use the package
The packages' GetCreds() function takes two string arguments. The first is the path to the Vault server. This defaults to `http://localhost:8200`, as this package is intended to use in a k8s environment with a Vault Agent sidecar loaded.
The second argument is the path to the where the SP will be created in Vault. This path is usually "/azure/creds/", followed by the name of the app, e.g "plattform-github-runner-scaler-test".

### Example
```
import vault "github.com/nrkno/plattform-vault-azure-sp-auth"

var spCreds *models.AzureCredentials
spCreds, err = vault.GetCreds("", VAULT_AZURE_SERVICE_PRINICIPAL_PATH)
if err != nil {
	log.Fatal("ERROR: Unable to get Azure SP credentials: " + err.Error())
}
```  
### Authenticating using the Service Principal
The GetCreds() function returns an AzureCredentials struct, consisting of two strings, ClientId and ClientSecret, which can be used for authentication when requesting Azure resources. The Service Principal created is valid for 24 hours, and must be recreated after expiry.  

## How to set up a Vault Agent sidecar in k8s
This package is intended to be used in a k8s environment, where a Vault Agent sidecar is loaded, and will forward requests to the Vault server. In order to set up a Vault Agent sidecar, the following annotations should be used:
```
vault.hashicorp.com/agent-cache-enable: "true"
vault.hashicorp.com/agent-cache-use-auto-auth-token: "true"
vault.hashicorp.com/agent-inject: "true"
vault.hashicorp.com/agent-init-first: "true"
vault.hashicorp.com/role: "your-role-here"
```
You need to have this *not* be set, as it the vault agent needs to be running in the pod: `vault.hashicorp.com/agent-pre-populate-only: "true"`