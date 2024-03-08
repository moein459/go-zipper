package main

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/moein459/go-zipper/api"
)

var validate *validator.Validate

func main() {
	engine := html.New("./web/views", ".html")
	validate = validator.New()
	validate.RegisterValidation("fileName", func(fl validator.FieldLevel) bool {
		match, _ := regexp.MatchString("^.*\\.(zip|ZIP)$", fl.Field().String())
		return match
	})

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(api.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		p := new(api.GenerateZipRequest)

		if err := c.BodyParser(p); err != nil {
			return c.Render("index", fiber.Map{"Error": err})
		}

		if err := validate.Struct(p); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return c.Render("index", fiber.Map{"ValidationErrors": validationErrors})
		}

		fmt.Println(p.FileName)
		return c.SendString("Validation is OK!")
	})

	app.Listen(":3000")
}
