package main

import (
	"github.com/gin-gonic/gin"
	"summercourse2/controller"
	"summercourse2/model"
	"summercourse2/service"
)

func main() {
	model.InitDB()
	service.InitService()

	r := gin.Default()
	r.GET("/getGoods", controller.SelectGoods)
	r.POST("/order", controller.MakeOrder)
	r.POST("/addGoods", controller.Addgoods)

	r.Run(":8080")
}
