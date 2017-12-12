package web

import (
	"fmt"
	"github.com/euclid1990/go-bigquery/web/controllers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"os"
)

func Init() {
	port := fmt.Sprintf(":%v", os.Getenv("WEB_PORT"))
	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New())

	// Load all templates from the "./web/views" folder inside cli context
	tmpl := iris.HTML("./web/views", ".html")
	tmpl.Layout("layout.html")
	tmpl.Reload(true)
	app.RegisterView(tmpl)

	// Define routing
	app.Controller("/", new(controllers.AppController))

	// Favicon & Static files serve
	app.Favicon("./web/assets/images/favicon.png")
	app.StaticWeb("/static", "./web/assets")

	app.Run(iris.Addr(port), iris.WithoutVersionChecker)
}
