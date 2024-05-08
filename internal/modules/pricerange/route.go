package pricerange

import (
	"fmt"
	"peluang-server/domain"
	"peluang-server/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type route struct {
	priceRangeService domain.PriceRangeService
}

func NewRoute(app *fiber.App, priceRangeService domain.PriceRangeService) {
	route := route{
		priceRangeService,
	}

	api := app.Group("/api")
	{
		api.Get("/pricerange", route.GetAllPriceRange)
		api.Post("/pricerange", route.CreatePriceRange)
	}
}

func (r *route) GetAllPriceRange(c *fiber.Ctx) error {
	pricerangeRes, err := r.priceRangeService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&dto.HttpResponse{
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
				Data:    []string{},
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		&dto.HttpResponse{
			Message: "success",
			Code:    fiber.StatusOK,
			Data:    pricerangeRes,
		},
	)
}

func (r *route) CreatePriceRange(c *fiber.Ctx) error {
	pricerange := new(dto.PriceRangeRequest)

	if err := c.BodyParser(pricerange); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&dto.HttpResponse{
				Message: "error parsing body",
				Code:    fiber.StatusBadRequest,
				Data:    []string{},
			},
		)
	}

	if err := validator.New().Struct(pricerange); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&dto.HttpResponse{
				Message: err.Error(),
				Code:    fiber.StatusBadRequest,
				Data:    []string{},
			},
		)
	}

	priceRangeModel := new(domain.PriceRange)
	priceRangeModel.Min = pricerange.Min
	priceRangeModel.Max = pricerange.Max

	slug := fmt.Sprintf("%d - %d", priceRangeModel.Min, priceRangeModel.Max)
	priceRangeModel.Slug = slug

	pricerangeRes, err := r.priceRangeService.Create(priceRangeModel, c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&dto.HttpResponse{
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
				Data:    []string{},
			},
		)
	}

	res := dto.PriceRangeResponse{
		Min:  pricerangeRes.Min,
		Max:  pricerangeRes.Max,
		Slug: pricerangeRes.Slug,
	}

	return c.Status(fiber.StatusCreated).JSON(
		&dto.HttpResponse{
			Message: "successfully created the price range",
			Code:    fiber.StatusCreated,
			Data:    res,
		},
	)
}
