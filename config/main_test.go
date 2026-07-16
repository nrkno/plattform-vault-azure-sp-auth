package config

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nrkno/plattform-vault-azure-sp-auth/models"
	"github.com/nrkno/plattform-vault-azure-sp-auth/utils"
)

// TestReadVaultPath_NilSecretReturnsError verifies that a nil secret returned by
// Vault (e.g. path not found, HTTP 404) results in a descriptive error rather
// than a nil pointer dereference panic.
func TestReadVaultPath_NilSecretReturnsError(t *testing.T) {
	// The Vault SDK returns (nil, nil) when the server responds with 404.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	retryCount := 0
	vaultCfg := models.VaultConfig{
		VaultAddress:                   server.URL,
		VaultAzureRolesCredentialsPath: "secret/nonexistent",
	}

	_, err := ReadVaultPath[models.AzureCredentials](vaultCfg, &ReadVaultPathOptions{
		RetryCount: utils.ToPointer(retryCount),
	})

	if err == nil {
		t.Fatal("expected error when vault returns no secret, got nil")
	}
	if !strings.Contains(err.Error(), "returned no secret") {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}
