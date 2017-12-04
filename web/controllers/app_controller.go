package controllers

import (
	"github.com/kataras/iris/mvc"
)

type AppController struct {
	mvc.C
}

func (c *AppController) Get() mvc.View {
	return mvc.View{
		Name: "app/index.html",
		Data: map[string]interface{}{
			"Title": "BigQuery & Golang App",
		},
	}
}
