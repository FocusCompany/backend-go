package api

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

var (
	jsonMarshal = jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
	}
)

func Init() {
	router := routing.New()
	router.Post("/window", RequireBasicJwt, windowHandler)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}
