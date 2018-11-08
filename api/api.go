package api

import (
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

func fooHandler(request *routing.Context) error {
	userId, err := ValidateJwtFromRequest(request)
	if err != nil {
		return err
	}

	fmt.Fprintf(request, "Hello, world, userID", userId.String())
	return nil
}

func Init() {
	router := routing.New()

	router.Get("/foo", fooHandler)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}
