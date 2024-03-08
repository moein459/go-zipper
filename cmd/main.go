package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
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

		rootFileDir := "files/"

		requestId := uuid.New().String()
		workingDir := filepath.Join(rootFileDir, requestId)

		os.Mkdir(workingDir, os.ModePerm)

		contentFileName := strings.TrimSuffix(p.FileName, filepath.Ext(p.FileName)) + ".txt"
		contentFile, _ := os.Create(filepath.Join(workingDir, contentFileName))
		contentFile.Write([]byte(p.Content))
		defer contentFile.Close()

		zipFileAddress := filepath.Join(workingDir, p.FileName)
		cmd := exec.Command("zip", "-j", "--password", p.Password, zipFileAddress, contentFile.Name())
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		return c.Download(zipFileAddress)
	})

	app.Listen(":3000")
}
