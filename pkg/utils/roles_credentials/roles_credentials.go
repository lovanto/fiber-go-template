package roles_credentials

import (
	"fmt"

	"github.com/create-go-app/fiber-go-template/pkg/repository"
)

var VerifyRole = func(role string) (string, error) {
	switch role {
	case repository.AdminRoleName:
	case repository.ModeratorRoleName:
	case repository.UserRoleName:
		// Nothing to do
	default:
		return "", fmt.Errorf("role '%v' does not exist", role)
	}
	return role, nil
}

func GetCredentialsByRole(role string) ([]string, error) {
	var credentials []string
	switch role {
	case repository.AdminRoleName:
		credentials = []string{
			repository.BookCreateCredential,
			repository.BookUpdateCredential,
			repository.BookDeleteCredential,
		}
	case repository.ModeratorRoleName:
		credentials = []string{
			repository.BookCreateCredential,
			repository.BookUpdateCredential,
		}
	case repository.UserRoleName:
		credentials = []string{
			repository.BookCreateCredential,
		}
	default:
		return nil, fmt.Errorf("role '%v' does not exist", role)
	}

	return credentials, nil
}
