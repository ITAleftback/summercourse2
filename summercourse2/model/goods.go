package model

import "github.com/jinzhu/gorm"
//定义一个商品结构体
type Goods struct {
	gorm.Model
	Name  string//商品名称
	Price int//商品价格
	Num   int//商品库存
}

// 添加商品的数据库操作
func (goods *Goods)AddGoods() error{
	return DB.Create(goods).Error
}

// 查看商品
func SelectGoodsById(id uint) (goods Goods, err error){
	err = DB.Table("goods").Where("id = ?",id).First(&goods).Error
	if err != nil {
		return Goods{}, err
	}
	return goods, nil
}

// 查看所有的商品
func SelectGoods() (goods []Goods, err error){
	err = DB.Table("goods").Find(&goods).Error
	if err != nil {
		return nil, err
	}
	return goods, nil
}
