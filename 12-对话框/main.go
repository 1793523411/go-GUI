package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-gtk/gtk"
)

func main() {
	gtk.Init(&os.Args)

	builder := gtk.NewBuilder()       //新建builder
	builder.AddFromFile("test.glade") //读取glade文件

	// 获取相应控件
	window := gtk.WindowFromObject(builder.GetObject("window1"))
	button1 := gtk.ButtonFromObject(builder.GetObject("button1"))
	button2 := gtk.ButtonFromObject(builder.GetObject("button2"))

	//信号处理
	button1.Connect("clicked", func() {
		//新建消息对话框，选择对话框
		dialog := gtk.NewMessageDialog(
			button1.GetTopLevelAsWindow(), //指定父窗口
			gtk.DIALOG_MODAL,              //模态对话框
			gtk.MESSAGE_QUESTION,          //指定对话框类型
			gtk.BUTTONS_YES_NO,            //默认按钮
			"Are u ok?")                   //设置内容

		dialog.SetTitle("问题对话框") //对话框设置标题

		flag := dialog.Run() //运行对话框
		if flag == gtk.RESPONSE_YES {
			fmt.Println("按下yes")
		} else if flag == gtk.RESPONSE_NO {
			fmt.Println("按下no")
		} else {
			fmt.Println("按下关闭按钮")
		}

		dialog.Destroy() //销毁对话框
	})

	button2.Connect("clicked", func() {
		dialog := gtk.NewMessageDialog(
			window,           //指定父窗口
			gtk.DIALOG_MODAL, //模态对话框
			gtk.MESSAGE_INFO, //info类型
			gtk.BUTTONS_OK,   //默认按钮
			"结束了")            //设置内容

		dialog.Run()     //运行对话框
		dialog.Destroy() //销毁对话框
	})

	window.Connect("destroy", gtk.MainQuit) //关闭窗口

	window.ShowAll()

	gtk.Main()
}
