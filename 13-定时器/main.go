package main

import (
	"os"
	"strconv"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	builder := gtk.NewBuilder()       //新建builder
	builder.AddFromFile("test.glade") //读取glade文件

	// 获取相应控件
	window := gtk.WindowFromObject(builder.GetObject("window1"))
	label := gtk.LabelFromObject(builder.GetObject("label"))
	buttonStart := gtk.ButtonFromObject(builder.GetObject("buttonStart"))
	buttonStop := gtk.ButtonFromObject(builder.GetObject("buttonStop"))

	label.ModifyFontSize(50)       //设置label字体大小
	buttonStop.SetSensitive(false) //停止按钮不能按

	var id int      //定时器id
	var num int = 1 //累加标记

	//信号处理
	//启动按钮
	buttonStart.Connect("clicked", func() {
		//启动定时器, 500毫秒为时间间隔，回调函数为匿名函数
		id = glib.TimeoutAdd(500, func() bool {
			num++
			label.SetText(strconv.Itoa(num)) //给标签设置内容
			return true                      //只要定时器没有停止，时间到自动调用回调函数
		})

		buttonStart.SetSensitive(false) //启动按钮变灰，不能按
		buttonStop.SetSensitive(true)   //定时器启动后，停止按钮可以按
	})

	//停止按钮
	buttonStop.Connect("clicked", func() {
		//停止定时器
		glib.TimeoutRemove(id)

		buttonStart.SetSensitive(true)
		buttonStop.SetSensitive(false)
	})

	window.Connect("destroy", gtk.MainQuit) //关闭窗口

	window.ShowAll()

	gtk.Main()
}

//import "github.com/mattn/go-gtk/glib"

//func TimeoutAdd(interval uint, f interface{}, datas ...interface{}) (id int)
//功能：创建定时器
//参数：
//    interval：设置的时间间隔，以毫秒为单位(1000即为1秒 )
//    f：回调函数的名字，回调函数的返回类型为bool，当回调函数返回值为false时，定时器执行一次后便会停止工作，不再循环执行。所以，要想定时器连续工作，循环执行所指定的回调函数，应该返回true。
//    datas：给回调函数传的参数
//返回值：定时器id号

//func TimeoutRemove(id int)
//功能：移除定时器
//参数：定时器id号
