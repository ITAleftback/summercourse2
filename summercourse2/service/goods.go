package service

import (

	"log"
	"summercourse2/model"
)

type Goods struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Num   int    `json:"num"`
}

// 添加商品
func AddGoods(name string,price int,num int) {
	// TODO
	//给商品结构体赋值
	good:=model.Goods{
		Name:  name,
		Price: price,
		Num:   num,
	}
	//调用modle的数据库操作
	err := good.AddGoods()
	if err != nil {
		log.Printf("Error Add goods. Error: %s",err)
	}
	log.Println("success")

}
//查询所有订单
func SelectGoods() (goods []Goods) {
	//拿到所有的商品信息
	_goods, err := model.SelectGoods()
	if err != nil {
		log.Printf("Error get goods info. Error: %s", err)
	}
	//goods是商品的集合一一拿到每个商品
	for _, v := range _goods {
		good := Goods{
			ID:    v.ID,
			Name:  v.Name,
			Price: v.Price,
			Num:   v.Num,
		}

		goods = append(goods, good)
	}
	//返回good集合
	return goods
}
