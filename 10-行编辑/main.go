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

	// 获取窗口控件指针，注意"window1"要和glade里的标志名称匹配
	window := gtk.WindowFromObject(builder.GetObject("window1"))
	window.SetSizeRequest(480, 320) //设置窗口大小

	//获取entry控件
	entry := gtk.EntryFromObject(builder.GetObject("entry1"))

	entry.SetText("123456")                       //设置内容
	fmt.Println("entry text = ", entry.GetText()) //获取内容
	//entry.SetVisibility(false)                  //设置不可见字符，即密码模式
	//entry.SetEditable(false)                    //只读，不可编辑
	entry.ModifyFontSize(30) //修改字体大小

	//信号处理，当用户在文本输入控件内部按回车键时引发activate信号
	entry.Connect("activate", func() {
		fmt.Println("entry text = ", entry.GetText()) //获取内容
	})

	//按窗口关闭按钮，自动触发"destroy"信号
	window.Connect("destroy", gtk.MainQuit)

	window.ShowAll()

	gtk.Main()
}

//信号标识 	触发条件
//“activate” 	行编辑区内部按回车键时触发
