package controllers

import (
	"context"
	"errors"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/repository"
	"github.com/create-go-app/fiber-go-template/pkg/utils/jwt"
	"github.com/create-go-app/fiber-go-template/pkg/utils/roles_credentials"
	"github.com/create-go-app/fiber-go-template/pkg/utils/wrapper"
	"github.com/create-go-app/fiber-go-template/platform/cache"
	"github.com/create-go-app/fiber-go-template/platform/database"

	"github.com/gofiber/fiber/v2"
)

// RenewTokens method for renew access and refresh tokens.
// @Description Renew access and refresh tokens.
// @Summary renew access and refresh tokens
// @Tags Token
// @Accept json
// @Produce json
// @Param request body models.Renew true "Renew Request"
// @Success 200 {object} models.TokenResponse
// @Security ApiKeyAuth
// @Router /v1/token/renew [post]
func RenewTokens(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	expiresAccessToken := claims.Expires
	if now > expiresAccessToken {
		return wrapper.ErrorResponse(c, fiber.StatusUnauthorized, "", errors.New(repository.UnauthorizedErrorMessage))
	}

	renew := &models.Renew{}
	if err := c.BodyParser(renew); err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	expiresRefreshToken, err := jwt.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	if now < expiresRefreshToken {
		userID := claims.UserID
		db, err := database.OpenDBConnection()
		if err != nil {
			return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
		}

		foundedUser, err := db.GetUserByID(userID)
		if err != nil {
			return wrapper.ErrorResponse(c, fiber.StatusNotFound, "", errors.New(repository.NotFoundErrorMessage))
		}

		credentials, err := roles_credentials.GetCredentialsByRole(foundedUser.UserRole)
		if err != nil {
			return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", err)
		}

		tokens, err := jwt.GenerateNewTokens(userID.String(), credentials)
		if err != nil {
			return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
		}

		connRedis, err := cache.RedisConnection()
		if err != nil {
			return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
		}

		errRedis := connRedis.Set(context.Background(), userID.String(), tokens.Refresh, 0).Err()
		if errRedis != nil {
			return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", errRedis)
		}

		return wrapper.SuccessResponse(c, "", models.TokenResponse{
			Access:  tokens.Access,
			Refresh: tokens.Refresh,
		})
	} else {
		return wrapper.ErrorResponse(c, fiber.StatusUnauthorized, "", errors.New(repository.UnauthorizedErrorMessage))
	}
}
