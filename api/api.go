package api

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

func Init() {
	router := routing.New()

	api := router.Group("", SetCorsHeader)
	api.Options("*", SetCorsHeader)

	api.Post("/window", RequireBasicJwt, windowHandler)
	api.Post("/window/list", RequireBasicJwt, windowListHandler)
	api.Post("/process", RequireBasicJwt, processHandler)
	api.Post("/process/list", RequireBasicJwt, processesListHandler)

	api.Get("/filters", RequireBasicJwt, getFiltersHandler)
	api.Post("/filters", RequireBasicJwt, updateFiltersHandler)

	api.Get("/dnd", RequireBasicJwt, getTotalDndHandler)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}
