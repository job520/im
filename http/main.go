package main

import (
	"fmt"
	"im/http/middleware"
	"im/http/router"
)

func main() {
	r := router.NewRouter()
	r.Use(middleware.Test())
	if err := r.Run(":" + "8080"); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
