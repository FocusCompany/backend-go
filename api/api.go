package api

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

func Init() {
	router := routing.New()
	router.Post("/window", RequireBasicJwt, windowHandler)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}
