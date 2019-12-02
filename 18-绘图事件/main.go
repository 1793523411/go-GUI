package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	builder := gtk.NewBuilder()       //新建builder
	builder.AddFromFile("test.glade") //读取glade文件

	window := gtk.WindowFromObject(builder.GetObject("window1")) //获取窗口控件
	window.SetSizeRequest(480, 320)                              //设置窗口大小
	window.SetAppPaintable(true)                                 //允许窗口能绘图(重要)

	var w, h int //保存窗口的宽度和高度

	//改变窗口大小时，触发"configure-event"，然后手动刷新绘图区域，否则图片会重叠
	window.Connect("configure-event", func() {
		window.QueueDraw() //刷新绘图区域

		//获取窗口的宽度和高度
		window.GetSize(&w, &h)
		fmt.Println(w, h)
	})

	x := 0 //画笑脸起点

	//绘图（曝光）事件，其回调函数做绘图操作
	window.Connect("expose-event", func() {
		//指定窗口为绘图区域，在窗口上绘图
		painter := window.GetWindow().GetDrawable()
		gc := gdk.NewGC(painter)

		//设置背景图的pixbuf，其宽高和窗口一样，最后一个参数固定为false
		bk, _ := gdkpixbuf.NewPixbufFromFileAtScale("./image/bk.jpg", w, h, false)
		//设置笑脸的pixbuf，其宽高为50，最后一个参数固定为false
		face, _ := gdkpixbuf.NewPixbufFromFileAtScale("./image/face.png", 50, 50, false)

		//画图
		//bk：需要绘图的pixbuf，第5、6个参数为画图的起点（相对于窗口而言）
		//第3、4个参数习惯为0，第7、8个参数为-1则按pixbuf原来的尺寸绘图
		//gdk.RGB_DITHER_NONE不用抖动，最后2个参数习惯写0
		painter.DrawPixbuf(gc, bk, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
		painter.DrawPixbuf(gc, face, 0, 0, x, 100, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)

		//释放pixbuf资源
		bk.Unref()
		face.Unref()
	})

	button := gtk.ButtonFromObject(builder.GetObject("button1")) //获取窗口控件
	//按钮"clicked"信号处理
	button.Clicked(func() {
		x += 50 //增加笑脸起点坐标

		if x >= w { //如果超过窗口宽度，置0
			x = 0
		}

		window.QueueDraw() //刷新绘图区域
	})

	//按窗口关闭按钮，自动触发"destroy"信号
	window.Connect("destroy", gtk.MainQuit)

	window.ShowAll()

	gtk.Main()
}

//GTK界面只要有图片的地方，其底层实际上是通过绘图实现的。绘图实际上也是事件的一种，GTK中，绘图事件也叫曝光事件。绘图的操作需要放在事件回调函数里。
//1------ 绘图触发条件

//绘图时所触发的信号：expose-event。只要触发曝光事件信号”expose-event”，就会自动调用所连接的回调函数。

//这里需要注意的是，曝光事件信号 “expose-event” 默认的情况下，是自动触发的(当然也可以人为触发)，就算我们不作任何操作，”expose-event”
//信号也有可能自动触发。前面我们学习中，我们按一下按钮就人为触发 “clicked” 信号，按一下鼠标人为触发 “button-press-event” 信号，
//如果我们不操作按钮，不操作鼠标，其对应的信号永远不会触发。

//窗口状态（移动，初始化，按按钮……）改变，只用我们肉眼能看到窗口上有变化，它都会自动触发曝光事件信号”expose-event”，
//然后就自动会调用它所连接的回调函数，但是，它不是刷新窗口的全部区域，它是按需要局部刷新，哪儿变化了就刷新那个变化的区域。

//当然我们也可以人为触发曝光事件信号”expose-event”：

//func (v *Widget) QueueDraw()

//    1

//2------- 绘图API

//窗口默认不允许在其上面绘图，需要人为设置：

//    window.SetAppPaintable(true)

//绘图操作：

//    //指定窗口为绘图区域，在窗口上绘图
//    painter := window.GetWindow().GetDrawable()
//    gc := gdk.NewGC(painter)

//    //bk：绘图的pixbuf，需要提前指定宽度和高度
//    //x, y：画图的起点（相对于窗口而言）
//    //其它参数为固定写法
//    painter.DrawPixbuf(gc, bk, 0, 0, x, y, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
