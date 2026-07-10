package GetCreds

import (
	"errors"
	"testing"

	"github.com/nrkno/plattform-vault-azure-sp-auth/config"
	"github.com/nrkno/plattform-vault-azure-sp-auth/models"
)

// MockVaultReader implements VaultReader for testing
type MockVaultReader struct {
	ReadVaultPathFunc func(vaultCfg models.VaultConfig, opts *config.ReadVaultPathOptions) (*models.AzureCredentials, error)
}

func (m *MockVaultReader) ReadVaultPath(vaultCfg models.VaultConfig, opts *config.ReadVaultPathOptions) (*models.AzureCredentials, error) {
	if m.ReadVaultPathFunc != nil {
		return m.ReadVaultPathFunc(vaultCfg, opts)
	}
	return nil, errors.New("not implemented")
}

// TestGetCredsDefaults tests that default values are set correctly
func TestGetCredsDefaults(t *testing.T) {
	tests := []struct {
		name                           string
		vaultAddr                      string
		vaultAzureRolesCredentialsPath string
		expectedAddr                   string
		expectedPath                   string
	}{
		{
			name:                           "both empty - uses defaults",
			vaultAddr:                      "",
			vaultAzureRolesCredentialsPath: "",
			expectedAddr:                   "http://localhost:8200",
			expectedPath:                   "default_vault_azure_roles_path",
		},
		{
			name:                           "custom addr, empty path",
			vaultAddr:                      "https://vault.example.com",
			vaultAzureRolesCredentialsPath: "",
			expectedAddr:                   "https://vault.example.com",
			expectedPath:                   "default_vault_azure_roles_path",
		},
		{
			name:                           "empty addr, custom path",
			vaultAddr:                      "",
			vaultAzureRolesCredentialsPath: "custom/path",
			expectedAddr:                   "http://localhost:8200",
			expectedPath:                   "custom/path",
		},
		{
			name:                           "both custom - no defaults",
			vaultAddr:                      "https://vault.prod.com",
			vaultAzureRolesCredentialsPath: "secret/azure/creds",
			expectedAddr:                   "https://vault.prod.com",
			expectedPath:                   "secret/azure/creds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedCfg models.VaultConfig

			// Create a mock reader that captures the config
			mockReader := &MockVaultReader{
				ReadVaultPathFunc: func(vaultCfg models.VaultConfig, opts *config.ReadVaultPathOptions) (*models.AzureCredentials, error) {
					capturedCfg = vaultCfg
					return &models.AzureCredentials{
						ClientId:     "test-client-id",
						ClientSecret: "test-client-secret",
					}, nil
				},
			}

			// Call GetCreds with the test inputs
			creds, err := GetCreds(mockReader, tt.vaultAddr, tt.vaultAzureRolesCredentialsPath)

			// Verify no error
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify the config passed to ReadVaultPath has the expected defaults
			if capturedCfg.VaultAddress != tt.expectedAddr {
				t.Errorf("VaultAddress: got %q, want %q", capturedCfg.VaultAddress, tt.expectedAddr)
			}

			if capturedCfg.VaultAzureRolesCredentialsPath != tt.expectedPath {
				t.Errorf("VaultAzureRolesCredentialsPath: got %q, want %q", capturedCfg.VaultAzureRolesCredentialsPath, tt.expectedPath)
			}

			// Verify we got credentials back
			if creds == nil {
				t.Errorf("expected credentials, got nil")
			}

			if creds.ClientId != "test-client-id" {
				t.Errorf("ClientId: got %q, want %q", creds.ClientId, "test-client-id")
			}
		})
	}
}

// TestGetCredsError tests that errors from ReadVaultPath are properly returned
func TestGetCredsError(t *testing.T) {
	mockReader := &MockVaultReader{
		ReadVaultPathFunc: func(vaultCfg models.VaultConfig, opts *config.ReadVaultPathOptions) (*models.AzureCredentials, error) {
			return nil, errors.New("vault connection failed")
		},
	}

	creds, err := GetCreds(mockReader, "http://localhost:8200", "secret/path")

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if creds != nil {
		t.Errorf("expected nil credentials on error, got %v", creds)
	}

	if err.Error() != "vault connection failed" {
		t.Errorf("unexpected error message: got %q", err.Error())
	}
}

// TestGetCredsSuccess tests successful credential retrieval
func TestGetCredsSuccess(t *testing.T) {
	expectedCreds := &models.AzureCredentials{
		ClientId:     "my-client-id",
		ClientSecret: "my-client-secret",
	}

	mockReader := &MockVaultReader{
		ReadVaultPathFunc: func(vaultCfg models.VaultConfig, opts *config.ReadVaultPathOptions) (*models.AzureCredentials, error) {
			return expectedCreds, nil
		},
	}

	creds, err := GetCreds(mockReader, "http://localhost:8200", "secret/path")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if creds.ClientId != expectedCreds.ClientId {
		t.Errorf("ClientId: got %q, want %q", creds.ClientId, expectedCreds.ClientId)
	}

	if creds.ClientSecret != expectedCreds.ClientSecret {
		t.Errorf("ClientSecret: got %q, want %q", creds.ClientSecret, expectedCreds.ClientSecret)
	}
}
