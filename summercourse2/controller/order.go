package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"summercourse2/service"
)

//下单
func MakeOrder(ctx *gin.Context) {
	userId := ctx.PostForm("userId")
	goodsId := ctx.PostForm("goodsId")
	itemId,_ := strconv.Atoi(goodsId)
	//这里是为了完成异步操作所以用到通道防止阻塞
	service.OrderChan <- service.User{
		UserId:  userId,
		GoodsId: uint(itemId),
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"info": "success",
	})
}

