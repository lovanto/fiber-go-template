package roles_credentials

import (
	"testing"

	"github.com/create-go-app/fiber-go-template/pkg/repository"
)

func TestVerifyRole(t *testing.T) {
	validRoles := []string{repository.AdminRoleName, repository.ModeratorRoleName, repository.UserRoleName}
	for _, role := range validRoles {
		got, err := VerifyRole(role)
		if err != nil {
			t.Fatalf("expected no error for role %s, got %v", role, err)
		}
		if got != role {
			t.Fatalf("expected role %s, got %s", role, got)
		}
	}

	if _, err := VerifyRole("invalid"); err == nil {
		t.Fatalf("expected error for invalid role, got nil")
	}
}

func TestGetCredentialsByRole(t *testing.T) {
	tests := []struct {
		role        string
		want        []string
		expectError bool
	}{
		{
			role: repository.AdminRoleName,
			want: []string{
				repository.BookCreateCredential,
				repository.BookUpdateCredential,
				repository.BookDeleteCredential,
			},
		},
		{
			role: repository.ModeratorRoleName,
			want: []string{
				repository.BookCreateCredential,
				repository.BookUpdateCredential,
			},
		},
		{
			role: repository.UserRoleName,
			want: []string{
				repository.BookCreateCredential,
			},
		},
		{
			role:        "ghost",
			expectError: true,
		},
	}

	for _, tt := range tests {
		got, err := GetCredentialsByRole(tt.role)
		if tt.expectError {
			if err == nil {
				t.Fatalf("expected error for role %s, got nil", tt.role)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for role %s: %v", tt.role, err)
		}
		if len(got) != len(tt.want) {
			t.Fatalf("role %s: expected %d creds, got %d", tt.role, len(tt.want), len(got))
		}
		// naive comparison
		for i, cred := range tt.want {
			if got[i] != cred {
				t.Fatalf("role %s: expected credential %s at index %d, got %s", tt.role, cred, i, got[i])
			}
		}
	}
}
