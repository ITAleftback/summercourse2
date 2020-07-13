package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"summercourse2/service"
)
//查询商品的接口
func SelectGoods(ctx *gin.Context) {
	//这里是拿到返回的所有的商品
	goods := service.SelectGoods()
	//返回用户消息
	ctx.JSON(http.StatusOK, gin.H{
		"status": 200,
		"info": "success",
		"data": struct {
			Goods []service.Goods `json:"goods"`
		}{goods},
	})
}
//添加商品
func Addgoods(ctx *gin.Context){
	good_name:=ctx.PostForm("good_name")
	good_price:=ctx.PostForm("good_price")
	num:=ctx.PostForm("num")
	goodprice,_ := strconv.Atoi(good_price)
	number,_ := strconv.Atoi(num)
	//调用服务端的添加商品
	service.AddGoods(good_name,goodprice,number)
	//返回用户消息
	ctx.JSON(http.StatusOK,gin.H{
		"status":200,
		"info":"success",
	})

}
