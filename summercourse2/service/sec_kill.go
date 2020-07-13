package service

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"sync"
	"time"
)

type User struct {
	UserId string
	GoodsId  uint
}
//带缓冲的channel
var OrderChan = make(chan User, 1024)
//用map来更新item，我思考了下如果是直接在数据库更新会进行两次数据库操作，这样也许效率更高？又或者是应该用redis？
var ItemMap = make(map[uint]*Item)
//item与good的区别在于，good是存在数据库里的也就是摆在面上的，下单不会造成变动，
// 下单应该减的是map里面item的数据?所以这里是不是应该用redis更好呢？
type Item struct {
	ID        uint   // 商品id
	Name      string // 名字
	Total     int    // 商品总量
	Left      int    // 商品剩余数量
	IsSoldOut bool   // 是否售罄
	leftCh    chan int
	sellCh    chan int
	done      chan struct{}
	Lock      sync.Mutex
}

// TODO 写一个定时任务，每天定时从数据库加载数据到Map
func Update(){
	//开启定时任务
	go func() {
		//这里我设置的时长是一分钟一次，方便观察函数正确与否
		crontab := cron.New()
		_ = crontab.AddFunc("0 */1 * * * ?", UpdateItem)
		crontab.Start()
	}()
}
//将数据库里面的good信息更新到map里
func UpdateItem(){
	//拿到所有商品，有个SelectGoods函数可以直接拿到商品，我就不重新写一个函数拿了
	goods:=SelectGoods()
	for _,good:=range goods{
		item:=&Item{
			ID:        good.ID,
			Name:      good.Name,
			Total:     100,//总量就随便写个值把
			Left:      good.Num,//剩余数量就是good的库存
			IsSoldOut: false,

		}
		ItemMap[item.ID]=item
	}
	//这个用来看是否更新成功的，正式里面要注释掉
	for _,v:=range ItemMap{
		fmt.Println(v)
	}

}
//测试测试
func initMap() {
	item := &Item{
		ID:        1,
		Name:      "测试",
		Total:     100,
		Left:      100,
		IsSoldOut: false,
		leftCh:    make(chan int),
		sellCh:    make(chan int),
	}
	ItemMap[item.ID] = item
}
//根据itemID 拿到商品的相关信息，返回map
func getItem(itemId uint) *Item{
	return ItemMap[itemId]
}
//用来开启异步的协程，在initService函数里面开10个协程是否会阻塞
func order() {
	for {
		user := <- OrderChan
		item := getItem(user.GoodsId)
		item.SecKilling(user.UserId)
	}
}
//秒杀，上锁放在阻塞，毕竟秒杀是多个同时调用
func (item *Item) SecKilling(userId string) {

	item.Lock.Lock()
	defer item.Lock.Unlock()
	// 等价
	// var lock = make(chan struct{}, 1}
	// lock <- struct{}{}
	// defer func() {
	// 		<- lock
	// }
	if item.IsSoldOut {
		return
	}
	item.BuyGoods(1)

	MakeOrder(userId, item.ID,1)


}

// 定时下架
func (item *Item) OffShelve() {
	beginTime := time.Now()
	// 获取第二天时间
	//nextTime := beginTime.Add(time.Hour * 24)
	// 计算次日零点，即商品下架的时间
	//offShelveTime := time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), 0, 0, 0, 0, nextTime.Location())
	offShelveTime := beginTime.Add(time.Second*5)
	timer := time.NewTimer(offShelveTime.Sub(beginTime))

	<-timer.C
	delete(ItemMap, item.ID)
	close(item.done)

}
// 出售商品
func (item *Item) SalesGoods() {
	for {
		select {
		//如果已经卖出的数量大于item的数量就是卖完了
		case num := <-item.sellCh:
			if item.Left -= num; item.Left <= 0 {
				item.IsSoldOut = true
			}
			//库存的信息发送给带channel的库存？？
		case item.leftCh <- item.Left:

		case <-item.Done():
			log.Println("我自闭了")
			return
		}
	}
}
//关闭item的salegoods
func (item *Item) Done() <-chan struct{} {
	if item.done == nil {
		item.done = make(chan struct{})
	}
	d := item.done
	return d
}

func (item *Item) Monitor() {
	go item.SalesGoods()
}

// 获取剩余库存
func (item *Item) GetLeft() int {
	var left int
	left = <-item.leftCh
	return left
}

// 购买商品
func (item *Item) BuyGoods(num int) {
	item.sellCh <- num
}

//服务端的操作
func InitService() {
	initMap()
	Update()
	for _,item := range ItemMap{
		item.Monitor()
		go item.OffShelve()
	}

	for i := 0; i < 10; i++ {
		go order()

	}

}
