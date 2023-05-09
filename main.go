package main

import (
	"github.com/tttturtle-russ/ree/ree"
	"log"
	"net/http"
)

func main() {
	e := ree.New()
	e.GET("/hello", func(ctx *ree.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	e.POST("/post", func(ctx *ree.Context) {

	})
	log.Fatal(e.Start(":9091"))
}
