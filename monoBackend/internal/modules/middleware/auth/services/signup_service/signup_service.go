package signup_service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"said-and-dot-backend/internal/common/bcrypt"
	"said-and-dot-backend/internal/common/validator"
	"said-and-dot-backend/internal/database"
	auth_errors "said-and-dot-backend/internal/modules/middleware/auth/errors"
	"time"
)

type SignupInput struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

func (si SignupInput) Validate() []*validator.ValidationError {
	return validator.ValidateStruct(si)
}

type SignupService interface {
	Signup(ctx *fiber.Ctx) error
}

type signupService struct {
	db database.Database
}

func NewSignupService(db database.Database) SignupService {
	return signupService{db: db}
}

func (ss signupService) Signup(ctx *fiber.Ctx) error {
	var signupInput SignupInput

	if err := ctx.BodyParser(&signupInput); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if validationErrors := signupInput.Validate(); validationErrors != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(validationErrors)
	}

	if err := createNewUser(ss, signupInput); err != nil {
		switch {
		case errors.Is(err, auth_errors.ErrUserAlreadyExists):
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "User with that email / username already exists",
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err,
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully registered a new account",
	})
}

func createNewUser(ss signupService, si SignupInput) error {
	var isUserExists bool
	if err := ss.db.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE username = $1 OR email = $2)",
		si.Username, si.Email).Scan(&isUserExists); err != nil {
		return err
	} else if isUserExists {
		return auth_errors.ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.Hash(si.Password)
	if err != nil {
		return err
	}

	if _, err := ss.db.Exec(
		"INSERT INTO Users "+
			"(username, password_hash, first_name, last_name, email, created) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		si.Username, passwordHash, si.FirstName, si.LastName, si.Email, time.Now()); err != nil {
		return err
	}

	return nil
}
