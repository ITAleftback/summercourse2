

# 引言

为了方便阅读，我整理一下。

# 修改的代码

## 有关添加商品

```
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
```

```
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
```

```
// 添加商品的数据库操作
func (goods *Goods)AddGoods() error{
   return DB.Create(goods).Error
}
```

## 更新map

```
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
```

```
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
```

# 总结

## os

sync.Mutex、sync.RWMutex等有关sync的操作都是与上锁相关。个人认为用的最多和最重要的应是sync.Mutex及sync.Waitgroup我详细总结一下。

### sync.Mutex

假设现在有两个程序同时对某个全局变量修改，会给出哪个结果？下图有详细解释

![](C:\Users\Mechrevo\Pictures\QQ图片20200713141926.png)

**·任何时间段只允许一个goroutine在临界区运行**

**·未加锁的Mutex会引起panic**

**·公平，互斥锁相互排斥，谁抢到锁谁执行**

![20190223120524726 (1)](C:\Users\Mechrevo\Pictures\20190223120524726 (1).png)

#### sync.Mutex使用

```
mutex:=&syne.Mutex{}

mutex.Lock()
//Update共享变量（比如切片，结构体等）
mutex.Unlock()
```

### sync.Waitgroup

**·等待一组goroutine完成**

**·Add参数可以是负值；如果计数器小于0，panic**

**·当计数器为0时，阻塞在wait方法 的goroutine会被释放**

```
var wg sync.WaitGroup
wg.Add(10)
for i := 0; i < 10; i++ {
   go order()
   wg.Done()
}
wg.Wait()
```

如上段代码，wg.Add(10)添加了10，而每一次遍历for循环都会执行wg.Done()，经历10次后计数器为0，wg.Wait()前的goroutine被释放。



值得注意的是wg.Wait()可以多次使用

```
var wg sync.WaitGroup
wg.Add(10)
for i := 0; i < 10; i++ {
   go order()
   wg.Done()
}
wg.Wait()
wg.Wait()
```

而wg.Done()就不能多次了，因为这里又减了一次，计数器为-1，引起panic

```
var wg sync.WaitGroup
wg.Add(10)
for i := 0; i < 10; i++ {
   go order()
   wg.Done()
}
wg.Done()
wg.Wait()
```

## happens-before

happens-before是一个术语，通常定义如下:

假设A和B表示一个多线程的程序执行的两个操作。如果A happens-before B，那么A操作对内存的影响 将对执行B的线程(且执行B之前)可见。

无论使用哪种编程语言，有一点是相同的：如果操作A和B在相同的线程中执行，并且A操作的声明在B 之前，那么A happens-before B。

```
int A, B;
void foo() 
{  
// This store to A ...  
A = 5; 
// ... effectively becomes visible before the following loads. Duh!  
B = A * A;
}
```

值得一提的是，happens-before关系是可传递的：如果A happens-before B，同时B happens-before C，那么A happens-before C。

刚接触这个术语的人总是容易误解，这里必须澄清的是，happens-before并不是指时序关系，并不是 说A happens-before B就表示操作A在操作B之前发生。它就是一个术语，就像光年不是时间单位一 样。具体地说：
1. A happens-before B并不意味着A在B之前发生。
2. A在B之前发生并不意味着A happens-before B。

下面举例说明

**A happens-before B并不意味着A在B之前发生**

```
int A = 0; 
int B = 0; 
void main()
{ 
A = B + 1;// (1)  
B = 1; // (2) 
}
```

根据前面说明的规则，(1) happens-before (2)。但是，如果我们使用gcc -O2编译这个代码，编译器将 产生一些指令重排序。有可能执行顺序是这样子的：

```
将B的值取到寄存器
将B赋值为1 
将寄存器值加1后赋值给A
```

根据定义，操作(1)对内存的影响必须 在操作(2)执行之前对其可见。换句话说，对A的赋值必须有机会对B的赋值有影响。

但是在这个例子中，对A的赋值其实并没有对B的赋值有影响。即便(1)的影响真的可见，(2)的行为还是 一样。所以，这并不能算是违背happens-before规则。

**A在B之前发生并不意味着A happens-before B**

```
int isReady = 0; 
int answer = 0;
void publishMessage() 
{  
answer = 42; // (1)  
isReady = 1; // (2) 
} 
void consumeMessage() 
{ 
if (isReady) // (3) <-- Let's suppose this line reads 1 
printf("%d\n", answer); // (4)
}
```

根据程序的顺序，在(1)和(2)之间存在happens-before 关系，同时在(3)和(4)之间也存在happens-before关系。 除此之外，我们假设在运行时，isReady读到1(是由另一个线程在(2)中赋的值)。在这中情形下，我们可 知(2)一定在(3)之前发生。但是这并不意味着在(2)和(3)之间存在happens-before 关系!
happens-before 关系只在语言标准中定义的地方存在，这里并没有相关的规则说明(2)和(3)之间存在 happens-before关系，即便(3)读到了(2)赋的值。 

还有，由于(2)和(3)之间，(1)和(4)之间都不存在happens-before关系，那么(1)和(4)的内存交互也可能 被重排序 (要不然来自编译器的指令重排序，要不然来自处理器自身的内存重排序)。那样的话，即使(3) 读到1，(4)也会打印出“0“。

## 原子操作

```
var i64 uint64
atomic.AddUint64(&64,5)//增，第一个参数为指针，第二个是增的值
atomic.LoadUint64(&i64)//返回指针指向的值
atomic.CompareAndSwapUint64(&i64,5,50)//比较并交换，第一个参数是需要替换的指针，第二个为旧值，第三个为新值，返回bool类型
atomic.SwapUint64(&i64,5) //将第二个参数替换指针指向的旧值，并将旧值返回
atomic.StoreUint64(&i64,5) //函数会把值赋到指针中

```

## Go关于同步的规则

关于channel的happens-before在Go的内存模型中提到了三种情况：

1. 对一个channel的发送操作 happens-before 相应channel的接收操作完成
2.  关闭一个channel happens-before 从该Channel接收到最后的返回值0 
3. 不带缓冲的channel的接收操作 happens-before 相应channel的发送操作完成

## 关于Go并发

1. 通过golang中的 goroutine 与sync.Mutex进行，并发同步。看起来协程似乎是同一时间进行，但在微观上，实则在同一时间只有一个协程，协程交替运行，达到了宏观上的同一时间效果。
2. goroutine之间通过 channel进行通信,channel是和类型相关的 可以理解为channel是一种类型安全的管道。
3. 带缓冲的channel与不带缓冲channel区别：当生产者channel加入值后，消费者才会取出，如果消费者取出后则会阻塞因为生产者没有还没有存入数据，而生产者赋值后也会阻塞直到消费者取出。而带缓冲的channel，生产者可以一直存入数据直到达到设置的容量前都不会阻塞，消费者就可以一直取出数据而不需要重新等待生产者存入新的数据。
4. 用了channel的函数一定要开协程，否则会阻塞
