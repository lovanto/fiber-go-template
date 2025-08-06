package utils

import (
	"fmt"

	"github.com/create-go-app/fiber-go-template/pkg/repository"
)

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
