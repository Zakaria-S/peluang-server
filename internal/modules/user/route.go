package user

import (
	"fmt"
	"peluang-server/domain"
	"peluang-server/dto"

	//"peluang-server/internal/middleware"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type route struct {
	userService domain.UserService
}

func NewRoute(app *fiber.App, userService domain.UserService) {
	route := route{
		userService,
	}

	api := app.Group("/api")
	{
		api.Post("/auth/register", route.UserRegister)
		api.Post("/auth/login", route.UserLogin)
		api.Post("/auth/otp", route.ValidateOTP)
	}
	/*
		protectedApi := app.Group("/api")
		protectedApi.Use(middleware.Authenticate())
		{
			api.Get("/users", route.GetUser)

		}
	*/
	limitedApi := app.Group("/api")
	limitedApi.Use(limiter.New(limiter.Config{
		Max:        1,
		Expiration: time.Minute * 2,
	}))
	{
		limitedApi.Post("/auth/resend-otp", route.ResendOTP)
	}

}
func (r *route) UserRegister(c *fiber.Ctx) error {
	user := new(dto.RegisterRequest)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&dto.HttpResponse{
				Message: "error parsing body",
				Code:    fiber.StatusBadRequest,
				Data:    []string{},
			},
		)
	}

	if err := validator.New().Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&dto.HttpResponse{
				Message: err.Error(),
				Code:    fiber.StatusBadRequest,
				Data:    []string{},
			},
		)
	}

	userModel := new(domain.User)
	userModel.Email = user.Email
	userModel.Password = user.Password
	userModel.Username = user.Username
	userModel.Telp = user.Telp

	userRes, otp, err := r.userService.Register(userModel, c.Context())
	if err != nil {
		if err == domain.ErrEmailExist {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": domain.ErrEmailExist.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	res := dto.RegisterResponse{
		ID:       userRes.ID.String(),
		Username: userRes.Username,
		Email:    userRes.Email,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"OTP":  otp,
		"data": res,
	})
}

func (r *route) UserLogin(c *fiber.Ctx) error {
	userReq := new(dto.LoginRequest)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}

	if err := validator.New().Struct(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userModel := new(domain.User)
	userModel.Email = userReq.Email
	userModel.Password = userReq.Password
	if token, err := r.userService.Login(userModel, c.Context()); err != nil {

		if err == domain.ErrInvalidCredential {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": domain.ErrInvalidCredential.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	}
}

func (r *route) GetUser(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			&dto.HttpResponse{
				Message: "invalid token",
				Code:    fiber.StatusUnauthorized,
				Data:    []string{},
			},
		)
	}

	tokenString = parts[1]

	user, err := r.userService.GetUser(tokenString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&dto.HttpResponse{
				Message: fmt.Sprintf("error getting user: %v", err),
				Code:    fiber.StatusInternalServerError,
				Data:    []string{},
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		&dto.HttpResponse{
			Message: "success",
			Code:    fiber.StatusOK,
			Data:    user,
		},
	)
}

func (r *route) ValidateOTP(c *fiber.Ctx) error {
	otp := new(dto.OTPRequest)
	if err := c.BodyParser(otp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("error parsing body: %v", err),
		})
	}

	if err := validator.New().Struct(otp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err := r.userService.ValidateOTP(otp.UserID, otp.OTP)
	if err != nil {
		if err == domain.ErrInvalidOTP {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": domain.ErrInvalidOTP.Error(),
			})
		}
		if err == domain.ErrExpiredOTP {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": domain.ErrExpiredOTP.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("error validating otp: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "otp valid",
	})
}

func (r *route) ResendOTP(c *fiber.Ctx) error {
	otp := new(dto.ReOTPRequest)
	if err := c.BodyParser(otp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("error parsing body: %v", err),
		})
	}

	if err := validator.New().Struct(otp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	newOtp, err := r.userService.ResendOTP(otp.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("error resending otp: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("otp %d has been sent", newOtp),
		"otp":     newOtp,
	})
}
