package controllers

import (
	"context"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/create-go-app/fiber-go-template/platform/cache"
	"github.com/create-go-app/fiber-go-template/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserSignUp method to create a new user.
// @Description Create a new user.
// @Summary create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.SignUp true "Sign Up Request"
// @Success 200 {object} models.SignUp
// @Router /v1/user/sign/up [post]
func UserSignUp(c *fiber.Ctx) error {
	signUp := &models.SignUp{}

	if err := c.BodyParser(signUp); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	validate := utils.NewValidator()
	if err := validate.Struct(signUp); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	role, err := utils.VerifyRole(signUp.UserRole)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	user := &models.SignUp{}
	user.Email = signUp.Email
	user.Password = utils.GeneratePassword(signUp.Password)
	user.UserRole = role

	userCreate := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        user.Email,
		PasswordHash: utils.GeneratePassword(user.Password),
		UserStatus:   1, // 0 == blocked, 1 == active
		UserRole:     user.UserRole,
	}

	if err := validate.Struct(userCreate); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	if err := db.CreateUser(userCreate); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	return utils.SuccessResponse(c, "", user)
}

// UserSignIn method to auth user and return access and refresh tokens.
// @Description Auth user and return access and refresh token.
// @Summary auth user and return access and refresh token
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.SignIn true "Sign In Request"
// @Success 200 {object} models.TokenResponse
// @Router /v1/user/sign/in [post]
func UserSignIn(c *fiber.Ctx) error {
	signIn := &models.SignIn{}

	if err := c.BodyParser(signIn); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	foundedUser, err := db.GetUserByEmail(signIn.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "", err)
	}

	compareUserPassword := utils.ComparePasswords(foundedUser.PasswordHash, signIn.Password)
	if !compareUserPassword {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	credentials, err := utils.GetCredentialsByRole(foundedUser.UserRole)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	tokens, err := utils.GenerateNewTokens(foundedUser.ID.String(), credentials)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	userID := foundedUser.ID.String()

	connRedis, err := cache.RedisConnection()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	errSaveToRedis := connRedis.Set(context.Background(), userID, tokens.Refresh, 0).Err()
	if errSaveToRedis != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", errSaveToRedis)
	}

	return utils.SuccessResponse(c, "", models.TokenResponse{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	})
}

// UserSignOut method to de-authorize user and delete refresh token from Redis.
// @Description De-authorize user and delete refresh token from Redis.
// @Summary de-authorize user and delete refresh token from Redis
// @Tags User
// @Accept json
// @Produce json
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/user/sign/out [post]
func UserSignOut(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	userID := claims.UserID.String()

	connRedis, err := cache.RedisConnection()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	errDelFromRedis := connRedis.Del(context.Background(), userID).Err()
	if errDelFromRedis != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "", errDelFromRedis)
	}

	return utils.SuccessResponse(c, "", "ok")
}
