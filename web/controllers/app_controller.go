package controllers

import (
	"github.com/euclid1990/go-bigquery/utilities"
	"github.com/kataras/iris/mvc"
	"golang.org/x/net/context"
)

type AppController struct {
	mvc.C
}

var (
	ctx      = context.Background()
	bigquery *utilities.BigQuery
)

func (c *AppController) Get() mvc.View {
	return mvc.View{
		Name: "app/index.html",
		Data: map[string]interface{}{
			"Title": "BigQuery & Golang App",
		},
	}
}

func (c *AppController) GetUser() mvc.View {
	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":  "User Analytics",
			"Header": "User",
		},
	}
}

func (c *AppController) GetAccess() mvc.View {
	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":  "User Analytics",
			"Header": "Access",
		},
	}
}

func (c *AppController) GetRetention() mvc.View {
	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":  "User Analytics",
			"Header": "RetentionRate",
		},
	}
}

func (c *AppController) GetDevice() mvc.View {
	return mvc.View{
		Name: "app/user.html",
		Data: map[string]interface{}{
			"Title":  "User Analytics",
			"Header": "Device",
		},
	}
}
