package service

import (
	"log"
	"summercourse2/model"
)

// order，下单，这里拿到想要下单的用户id，商品id，数量并赋值给结构体后进行数据库操作
func MakeOrder(userId string, goodsId uint, num int) {

	order := model.Order{
		UserID:  userId,
		GoodsID: goodsId,
		Num:     num,
	}
	err := order.MakeOrder()
	if err != nil {
		log.Printf("Error make an order. Error: %s",err)
	}
	log.Println("success")
}