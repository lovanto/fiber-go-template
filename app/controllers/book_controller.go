package controllers

import (
	"errors"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/repository"
	"github.com/create-go-app/fiber-go-template/pkg/utils/jwt"
	"github.com/create-go-app/fiber-go-template/pkg/utils/validator"
	"github.com/create-go-app/fiber-go-template/pkg/utils/wrapper"
	"github.com/create-go-app/fiber-go-template/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetBooks func gets all exists books.
// @Description Get all exists books.
// @Summary get all exists books
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Router /v1/books [get]
func GetBooks(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	books, err := db.GetBooks()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "books were not found",
			"count": 0,
			"books": nil,
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(books),
		"books": books,
	})
}

// GetBook func gets book by given ID or 404 error.
// @Description Get book by given ID.
// @Summary get book by given ID
// @Tags Book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Router /v1/book/{id} [get]
func GetBook(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	book, err := db.GetBook(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "book with the given ID is not found",
			"book":  nil,
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"book":  book,
	})
}

// CreateBook func for creates a new book.
// @Description Create a new book.
// @Summary create a new book
// @Tags Book
// @Accept json
// @Produce json
// @Param request body models.BookCreate true "Book"
// @Success 201 {object} models.Book
// @Security ApiKeyAuth
// @Router /v1/book [post]
func CreateBook(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	expires := claims.Expires
	if now > expires {
		return wrapper.ErrorResponse(c, fiber.StatusUnauthorized, "", errors.New(repository.UnauthorizedErrorMessage))
	}

	credential := claims.Credentials[repository.BookCreateCredential]
	if !credential {
		return wrapper.ErrorResponse(c, fiber.StatusForbidden, "", errors.New(repository.ForbiddenErrorMessage))
	}

	book := &models.BookCreate{}

	if err := c.BodyParser(book); err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	validate := validator.NewValidator()

	bookCreate := &models.Book{}
	bookCreate.ID = uuid.New()
	bookCreate.CreatedAt = time.Now()
	bookCreate.UpdatedAt = time.Now()
	bookCreate.UserID = claims.UserID
	bookCreate.Title = book.Title
	bookCreate.Author = book.Author
	bookCreate.BookAttrs = book.BookAttrs
	bookCreate.BookStatus = 1 // 0 == draft, 1 == active

	if err := validate.Struct(book); err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", errors.New(validator.ValidatorErrors(err)))
	}

	if err := db.CreateBook(bookCreate); err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	return wrapper.SuccessResponse(c, "", book)
}

// UpdateBook func for updates book by given ID.
// @Description Update book.
// @Summary update book
// @Tags Book
// @Accept json
// @Produce json
// @Param request body models.BookUpdate true "Book"
// @Success 201 {object} models.Book
// @Security ApiKeyAuth
// @Router /v1/book [put]
func UpdateBook(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	expires := claims.Expires
	if now > expires {
		return wrapper.ErrorResponse(c, fiber.StatusUnauthorized, "", errors.New(repository.UnauthorizedErrorMessage))
	}

	credential := claims.Credentials[repository.BookUpdateCredential]
	if !credential {
		return wrapper.ErrorResponse(c, fiber.StatusForbidden, "", errors.New(repository.ForbiddenErrorMessage))
	}

	book := &models.BookUpdate{}
	if err := c.BodyParser(book); err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	foundedBook, err := db.GetBook(book.ID)
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusNotFound, "", errors.New(repository.NotFoundErrorMessage))
	}

	userID := claims.UserID

	if foundedBook.UserID == userID {
		bookUpdate := &models.Book{}
		bookUpdate.ID = book.ID
		bookUpdate.UpdatedAt = time.Now()
		bookUpdate.Title = book.Title
		bookUpdate.Author = book.Author
		bookUpdate.BookAttrs = book.BookAttrs
		validate := validator.NewValidator()
		if err := validate.Struct(book); err != nil {
			return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", errors.New(validator.ValidatorErrors(err)))
		}

		if err := db.UpdateBook(foundedBook.ID, bookUpdate); err != nil {
			return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
		}

		return wrapper.SuccessResponse(c, "", book)
	} else {
		return wrapper.ErrorResponse(c, fiber.StatusForbidden, "", errors.New(repository.ForbiddenDataModificationErrorMessage))
	}
}

// DeleteBook func for deletes book by given ID.
// @Description Delete book by given ID.
// @Summary delete book by given ID
// @Tags Book
// @Accept json
// @Produce json
// @Param request body models.BookDelete true "Book ID"
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/book [delete]
func DeleteBook(c *fiber.Ctx) error {
	now := time.Now().Unix()

	claims, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	expires := claims.Expires
	if now > expires {
		return wrapper.ErrorResponse(c, fiber.StatusUnauthorized, "", errors.New(repository.UnauthorizedErrorMessage))
	}

	credential := claims.Credentials[repository.BookDeleteCredential]
	if !credential {
		return wrapper.ErrorResponse(c, fiber.StatusForbidden, "", errors.New(repository.ForbiddenErrorMessage))
	}

	book := &models.BookDelete{}

	if err := c.BodyParser(book); err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	validate := validator.NewValidator()
	if err := validate.StructPartial(book, "id"); err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "", err)
	}

	db, err := database.OpenDBConnection()
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
	}

	foundedBook, err := db.GetBook(book.ID)
	if err != nil {
		return wrapper.ErrorResponse(c, fiber.StatusNotFound, "", errors.New(repository.NotFoundErrorMessage))
	}

	userID := claims.UserID
	if foundedBook.UserID == userID {
		if err := db.DeleteBook(foundedBook.ID); err != nil {
			return wrapper.ErrorResponse(c, fiber.StatusInternalServerError, "", err)
		}

		return wrapper.SuccessResponse(c, "", "ok")
	} else {
		return wrapper.ErrorResponse(c, fiber.StatusForbidden, "", errors.New(repository.ForbiddenDataModificationErrorMessage))
	}
}
