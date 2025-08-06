package utils

import (
	"fmt"

	"github.com/create-go-app/fiber-go-template/pkg/repository"
)

func VerifyRole(role string) (string, error) {
	switch role {
	case repository.AdminRoleName:
	case repository.ModeratorRoleName:
	case repository.UserRoleName:
		// Nothing to do, verified successfully.
	default:
		return "", fmt.Errorf("role '%v' does not exist", role)
	}

	return role, nil
}
