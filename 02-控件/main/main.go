package main

import (
	"os"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args) //环境初始化

	//--------------------------------------------------------
	// 主窗口
	//--------------------------------------------------------
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL) //创建窗口
	window.SetPosition(gtk.WIN_POS_CENTER)       //设置窗口居中显示
	window.SetTitle("GTK Go!")                   //设置标题
	window.SetSizeRequest(300, 200)              //设置窗口的宽度和高度

	//--------------------------------------------------------
	// GtkFixed
	//--------------------------------------------------------
	layout := gtk.NewFixed() //创建固定布局

	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	b1 := gtk.NewButton() //新建按钮
	b1.SetLabel("^_@")    //设置内容
	//b1.SetSizeRequest(100, 50) //设置按钮大小

	b2 := gtk.NewButtonWithLabel("@_~") //新建按钮，同时设置内容
	b2.SetSizeRequest(100, 50)          //设置按钮大小

	//--------------------------------------------------------
	// 添加布局、添加容器
	//--------------------------------------------------------
	window.Add(layout) //把布局添加到主窗口中

	layout.Put(b1, 0, 0)    //设置按钮在容器的位置
	layout.Move(b1, 50, 50) //移动按钮的位置，必须先put，再用move

	layout.Put(b2, 50, 100)

	window.ShowAll() //显示所有的控件

	gtk.Main() //主事件循环，等待用户操作
}

//func (v *Fixed) Put(w IWidget, x, y int)
//功能：固定布局容器添加控件
//参数：
//    widget：要添加的控件
//    x, y：控件摆放位置的起点坐标
