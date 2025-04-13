package main

import (
	"filepro/config"
	"filepro/model"
	"filepro/router"
	"filepro/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"log"
	"os"
)

func init() {
	utils.Createdir("data")
	utils.Createdir("zips")
	utils.Createdir("qr/pic")
}
func main() {
	config.InitConfig()
	os.Stat(model.Dir)
	os.MkdirAll(model.Dir, os.ModePerm)
	config.ConfigLock.Lock()
	app := fiber.New(fiber.Config{
		Views:     html.New("./template", ".html"),
		BodyLimit: 100000 * 1024 * 1024,
	})
	config.ConfigLock.Unlock()
	app.Static("/pic", "./qr/pic")
	app.Static("/static", "./template/static")

	index := app.Group("/")
	api := app.Group("/api")
	admin := app.Group("/admin")

	router.Router_page(index)
	router.Router_api(api)
	router.Router_admin(admin)

	app.Use(func(c *fiber.Ctx) error {
		log.Printf("WARN: Route not found - Path: %s", c.Path())
		c.Status(fiber.StatusNotFound)
		return c.Type("html").SendString(`
            <!DOCTYPE html>
            <html>
            <head><title>404 Not Found</title></head>
            <body>
                <h1>Oops! Page Not Found</h1>
                <p>Sorry, the resource you requested at path <code>` + c.Path() + `</code> could not be found.</p>
                <a href="/">Go to Homepage</a>
            </body>
            </html>
        `)
	})

	log.Println(app.Listen(":3000"))
}
