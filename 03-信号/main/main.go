package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

//按钮b1信号处理的回调函数
func HandleButton(ctx *glib.CallbackContext) {
	arg := ctx.Data()   //获取用户传递的参数，是空接口类型
	p, ok := arg.(*int) //类型断言
	if ok {             //如果ok为true，说明类型断言正确
		fmt.Println("*p = ", *p) //用户传递传递的参数为&tmp，是一个变量的地址
		*p = 250                 //操作指针所指向的内存
	}

	fmt.Println("按钮b1被按下")

	//gtk.MainQuit() //关闭gtk程序
}

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

	//--------------------------------------------------------
	// 信号处理
	//--------------------------------------------------------
	//按钮按下自动触发"pressed"，自动调用HandleButton, 同时将 &tmp 传递给HandleButton
	tmp := 10
	b1.Connect("pressed", HandleButton, &tmp)

	//回调函数为匿名函数，推荐写法
	//按钮按下自动触发"pressed"，自动调用匿名函数，
	b2.Connect("pressed", func() {

		fmt.Println("b2被按下")
		fmt.Println("tmp = ", tmp)

	}) //注意：}和)在同一行

	window.ShowAll() //显示所有的控件

	gtk.Main() //主事件循环，等待用户操作
}

//func (v *Widget) Connect(s string, f interface{}, datas ...interface{}) int
//功能：信号注册
//参数：
//    v: 信号发出者，可以认为我们操作的控件，如按下按钮，这个就为按钮指针
//    s：信号标志，如"pressed"
//    f：回调函数的名称，
//    datas：给回调函数传的参数，尽管是可变参数，但是只能传递一个参数，可变参数的目的为了让用户多个选择(可以传参，或者不传)
//返回值：
//    注册函数的标志
