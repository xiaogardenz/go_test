package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"golang.org/x/net/context"
	"log"
	"orange-iris/seckill_iris/backend/web/controllers"
	"orange-iris/seckill_iris/common"
	"orange-iris/seckill_iris/repositories"
	"orange-iris/seckill_iris/services"
)
//test分支加注释
func main() {
	//1.创建实例
	app := iris.New()
	//2。设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模版 Layout(布局文件)
	tmplate := iris.HTML("./backend/web/views",".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板目标
	app.HandleDir("/assets","./backend/web/assets")
	//5.出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message",ctx.Values().GetStringDefault("message","访问的页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//连接数据库
	db,err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
	}

	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	//6.注册控制器
	productRepository := repositories.NewProductManager("product",db)
	productSerivce := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx,productSerivce)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderMangerRepository("order",db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx,orderService)
	order.Handle(new(controllers.OrderController))

	/*作业
		1、继续完成订单管理的其他功能展示
			1）查询并且展示单个订单信息
			2）删除订单；
			3）更新修改订单
		2、简化controller 注册代码
			1）简单创建controller就实现注册
		3、如何简化model层代码（中小项目）
			1）推荐学习gorm
			2) 掌握gorm
	*/
	//7。启动服务
	app.Run(iris.Addr(":8080"))

}
