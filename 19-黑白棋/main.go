package main

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

//控件结构体
type ChessWidet struct {
	window      *gtk.Window //窗口
	buttonMin   *gtk.Button //最小化按钮
	buttonClose *gtk.Button //关闭按钮
	labelBlack  *gtk.Label  //显示黑子个数
	labelWhite  *gtk.Label  //显示白子个数
	labelTime   *gtk.Label  //倒计时间显示
	imageBlack  *gtk.Image  //黑子image，用于提示该谁落子
	imageWhite  *gtk.Image  //白子image，用于提示该谁落子
}

//属性结构体
type ChessInfo struct {
	w, h int //窗口的宽度高度
	x, y int //鼠标点击坐标(相对于窗口)

	gridW  int //棋盘水平方向一个格子的宽度
	gridH  int //棋盘水平方向一个格子的高度
	startX int //棋盘起点x坐标
	startY int //棋盘起点y坐标
}

//枚举，标志棋盘棋子状态
const (
	Empty = iota //当前棋盘格子没有子
	Black        //当前棋盘格子为黑子
	White        //当前棋盘格子为白子
)

type Chessboard struct {
	//匿名字段
	ChessWidet
	ChessInfo

	currentRole    int //当前落子角色(该谁落子)
	tipTimerId     int //定时器id，用于实现落子落子一闪一闪效果
	machineTimerId int //定时器id，机器1s后才落子
	leftTimerId    int //定时器id，用于倒计时
	timeNum        int //倒计时间

	chess [8][8]int //二维数组，标记棋盘棋子状态
}

//函数，给按钮设置图片
func ButtonSetImageFromFile(button *gtk.Button, filename string) {
	//获取按钮的宽度和高度
	var w, h int
	button.GetSizeRequest(&w, &h)
	//获取pixbuf，指定其大小
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale(filename, w-10, h-10, false)

	//通过pixbuf新建image
	image := gtk.NewImageFromPixbuf(pixbuf)
	//按钮设置图片
	button.SetImage(image)

	//释放pixbuf资源
	pixbuf.Unref()
}

//函数，给image设置图片
func ImageSetPixbufFromFile(image *gtk.Image, filename string) {
	//获取image的宽度和高度
	var w, h int
	image.GetSizeRequest(&w, &h)
	//获取pixbuf，指定其大小
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale(filename, w-5, h-5, false)

	//image设置pixbuf
	image.SetFromPixbuf(pixbuf)

	//释放pixbuf资源
	pixbuf.Unref()
}

//Chessboard的方法，获取控件，设置控件属性，返回主窗口控件指针
func (obj *Chessboard) CreateWindow() *gtk.Window {
	builder := gtk.NewBuilder()       //新建builder
	builder.AddFromFile("test.glade") //读取glade文件

	//获取glade对应的控件
	obj.window = gtk.WindowFromObject(builder.GetObject("window1"))      //获取窗口控件
	obj.buttonMin = gtk.ButtonFromObject(builder.GetObject("buttonMin")) //按钮
	obj.buttonClose = gtk.ButtonFromObject(builder.GetObject("buttonClose"))
	obj.labelBlack = gtk.LabelFromObject(builder.GetObject("labelBlack")) //标签
	obj.labelWhite = gtk.LabelFromObject(builder.GetObject("labelWhite"))
	obj.labelTime = gtk.LabelFromObject(builder.GetObject("labelTime"))
	obj.imageBlack = gtk.ImageFromObject(builder.GetObject("imageBlack")) //image
	obj.imageWhite = gtk.ImageFromObject(builder.GetObject("imageWhite")) //image

	//设置属性
	//窗口属性
	obj.w = 800
	obj.h = 480
	obj.window.SetSizeRequest(obj.w, obj.h)    //设置窗口大小
	obj.window.SetPosition(gtk.WIN_POS_CENTER) //居中显示
	obj.window.SetAppPaintable(true)           //允许窗口能绘图(重要)
	obj.window.SetDecorated(false)             //去表框
	//添加鼠标按下事件
	obj.window.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	//按钮属性
	ButtonSetImageFromFile(obj.buttonMin, "image/min.png") //给按钮设置图片，此为自定义函数
	ButtonSetImageFromFile(obj.buttonClose, "image/close.png")
	obj.buttonMin.SetCanFocus(false) //去掉按钮上的聚焦框
	obj.buttonClose.SetCanFocus(false)

	//标签属性
	obj.labelBlack.ModifyFontSize(50) //修改字体大小
	obj.labelWhite.ModifyFontSize(50) //修改字体大小
	obj.labelTime.ModifyFontSize(30)  //修改字体大小

	//修改字体颜色为白色
	obj.labelBlack.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	obj.labelWhite.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	obj.labelTime.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	//image属性
	ImageSetPixbufFromFile(obj.imageBlack, "image/black.png")
	ImageSetPixbufFromFile(obj.imageWhite, "image/white.png")

	//obj.imageBlack.Hide()
	//obj.imageWhite.Hide()

	//棋盘格子尺寸信息
	obj.startX = 200
	obj.startY = 60
	obj.gridW = 50
	obj.gridH = 40

	return obj.window
}

//绘图事件处理函数，"expose-event"的回调函数
func PaintEvent(ctx *glib.CallbackContext) {
	arg := ctx.Data()            //获取用户传递的参数，是空接口类型
	obj, ok := arg.(*Chessboard) //类型断言
	if false == ok {             //如果ok为false，说明类型断言错误
		fmt.Println("arg.(*Chessboard) err")
		return
	}

	//指定窗口为绘图区域，在窗口上绘图
	painter := obj.window.GetWindow().GetDrawable()
	gc := gdk.NewGC(painter)

	//设置背景图的pixbuf，其宽高和窗口一样，最后一个参数固定为false
	bg, _ := gdkpixbuf.NewPixbufFromFileAtScale("./image/bg.jpg", obj.w, obj.h, false)

	//画图，画背景图
	painter.DrawPixbuf(gc, bg, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)

	//设置黑白子图片pixbuf
	blackPixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale("./image/black.png", obj.gridW, obj.gridH, false)
	whitePixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale("./image/white.png", obj.gridW, obj.gridH, false)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if obj.chess[i][j] == Black { //画黑子
				painter.DrawPixbuf(gc, blackPixbuf, 0, 0, obj.startX+i*obj.gridW, obj.startY+j*obj.gridH, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
			} else if obj.chess[i][j] == White { //画黑子
				painter.DrawPixbuf(gc, whitePixbuf, 0, 0, obj.startX+i*obj.gridW, obj.startY+j*obj.gridH, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
			}
		}
	}

	//释放图片资源，必须，否则会导致内存泄露，内存越用越多
	bg.Unref()
	blackPixbuf.Unref()
	whitePixbuf.Unref()
}

//方法，统计黑白棋个数，胜负判断
func (obj *Chessboard) JudgeResult() {
	var blackNum, whiteNum int //黑白子个数
	isOver := true             //标记游戏是否结束，true为结束，false为没有结束

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if obj.chess[i][j] == Black { //黑子
				blackNum++ //黑子累加
			} else if obj.chess[i][j] == White {
				whiteNum++ //白子累加
			}

			//判断黑白子在i, j位置，能否吃子，没有改变chess的标志位
			if 0 < obj.JudgeRule(i, j, Black, false) || 0 < obj.JudgeRule(i, j, White, false) {
				isOver = false //只要能吃子，说明游戏还没有结束
			}

		}
	}

	//黑白棋个数显示
	obj.labelBlack.SetText(strconv.Itoa(blackNum))
	obj.labelWhite.SetText(strconv.Itoa(whiteNum))

	//fmt.Println("isOver = ", isOver)
	if isOver == false {
		return //游戏还没有结束，终止此函数
	}

	//如果游戏结束，中断定时器
	glib.TimeoutRemove(obj.tipTimerId)
	glib.TimeoutRemove(obj.leftTimerId)

	//判断胜负
	var str string
	if blackNum > whiteNum {
		str = "你(蓝蜘蛛)赢了！！！\n继续游戏，按\"是\""
	} else if blackNum < whiteNum {
		str = "机器(橙蜘蛛)赢了！！！\n继续游戏，按\"是\""
	} else {
		str = "平局！！！\n继续游戏，按\"是\""
	}

	//弹出对话框
	dialog := gtk.NewMessageDialog(
		obj.window,         //指定父窗口
		gtk.DIALOG_MODAL,   //模态对话框
		gtk.MESSAGE_INFO,   //info类型
		gtk.BUTTONS_YES_NO, //默认按钮
		str)                //设置内容
	dialog.SetTitle("游戏结束")

	result := dialog.Run() //运行对话框
	if result == gtk.RESPONSE_YES {
		fmt.Println("按下yes")
		obj.InitChess() //重新初始化数据，继续游戏

	} else {
		fmt.Println("按下关闭按钮")
		//关闭程序
		gtk.MainQuit()
	}

	dialog.Destroy() //销毁对话框
}

//函数，机器落子
func MachinePlay(obj *Chessboard) {
	//先暂停定时器
	glib.TimeoutRemove(obj.machineTimerId)

	max, px, py := 0, -1, -1

	//找到能吃子的位置
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			num := obj.JudgeRule(i, j, obj.currentRole, false) //判断能吃几个棋子
			if num > 0 {
				//先找4个角落，先占角
				if (i == 0 && j == 0) || (i == 7 && j == 0) || (i == 0 && j == 7) || (i == 7 && j == 7) {
					px, py = i, j
					goto End
				}

				//不能占角，找吃子最多的
				if max < num { //吃子个数最多的
					px, py, max = i, j, num
				}
			}
		}
	}

End:
	if px == -1 { //说明机器没有落子的地方
		obj.ChangeRole() //改变角色
		return
	}

	obj.JudgeRule(px, py, obj.currentRole, true) //机器吃子
	obj.window.QueueDraw()                       //刷新绘图区域
	obj.ChangeRole()                             //改变落子角色

}

//方法，改变落子角色
func (obj *Chessboard) ChangeRole() {
	//重新设置时间
	obj.timeNum = 20
	obj.labelTime.SetText(strconv.Itoa(obj.timeNum)) //标签显示时间

	//先隐藏提示图片
	obj.imageBlack.Hide()
	obj.imageWhite.Hide()

	if obj.currentRole == Black { //白子下
		obj.currentRole = White
	} else { //黑子下
		obj.currentRole = Black
	}

	//统计黑白棋个数，胜负判断
	obj.JudgeResult()

	if obj.currentRole == White { //机器落子，启动定时器
		obj.machineTimerId = glib.TimeoutAdd(1000, func() bool {
			MachinePlay(obj) //机器落子
			return true
		})
	}

}

// 吃子的规则
// 吃子规则的参数：棋盘数组坐标位置(x y) role 当前落子角色
// eatChess为true，代表改变原来的数组， false不改变数组内容，只判断此位置能吃多少个子
// 返回值：吃子个数
func (obj *Chessboard) JudgeRule(x, y int, role int, eatChess bool) (eatNum int) {
	// 棋盘的八个方向
	dir := [8][2]int{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}}

	tempX, tempY := x, y // 临时保存棋盘数组坐标位置

	if obj.chess[tempX][tempY] != Empty { // 如果此方格内已有棋子，返回
		return 0
	}

	// 棋盘的8个方向
	for i := 0; i < 8; i++ {
		tempX += dir[i][0]
		tempY += dir[i][1] // 准备判断相邻棋子

		// 如果没有出界，且相邻棋子是对方棋子，才有吃子的可能．
		if (tempX < 8 && tempX >= 0 && tempY < 8 && tempY >= 0) && (obj.chess[tempX][tempY] != role) && (obj.chess[tempX][tempY] != Empty) {
			tempX += dir[i][0]
			tempY += dir[i][1] // 继续判断下一个，向前走一步
			for tempX < 8 && tempX >= 0 && tempY < 8 && tempY >= 0 {
				if obj.chess[tempX][tempY] == Empty { // 遇到空位跳出
					break
				}

				if obj.chess[tempX][tempY] == role { // 找到自己的棋子，代表可以吃子
					if eatChess == true { // 确定吃子
						obj.chess[x][y] = role // 开始点标志为自己的棋子
						tempX -= dir[i][0]
						tempY -= dir[i][1] // 后退一步
						for (tempX != x) || (tempY != y) {
							// 只要没有回到开始的位置就执行
							obj.chess[tempX][tempY] = role // 标志为自己的棋子
							tempX -= dir[i][0]
							tempY -= dir[i][1] // 继续后退一步
							eatNum++           // 累计
						}
					} else { //不吃子，只是判断这个位置能不能吃子
						tempX -= dir[i][0]
						tempY -= dir[i][1]                 // 后退一步
						for (tempX != x) || (tempY != y) { // 只计算可以吃子的个数
							tempX -= dir[i][0]
							tempY -= dir[i][1] // 继续后退一步
							eatNum++
						}
					}

					break // 跳出循环
				} // 没有找到自己的棋子，就向前走一步

				tempX += dir[i][0]
				tempY += dir[i][1] // 向前走一步
			}
		} // 如果这个方向不能吃子，就换一个方向
		tempX, tempY = x, y
	}

	return // 返回能吃子的个数
}

//鼠标按下事件处理，MousePressEvent为其回调函数，把obj传递给回调函数
func MousePressEvent(ctx *glib.CallbackContext) {
	arg := ctx.Data()            //获取用户传递的参数，是空接口类型
	obj, ok := arg.(*Chessboard) //类型断言
	if false == ok {             //如果ok为false，说明类型断言错误
		fmt.Println("arg.(*Chessboard) err")
		return
	}

	//获取鼠键按下属性结构体变量，系统内部的变量，不是用户传参变量
	tmp := ctx.Args(0)
	event := *(**gdk.EventButton)(unsafe.Pointer(&tmp))
	if event.Button == 1 { //左键
		obj.x, obj.y = int(event.X), int(event.Y) //保存点击的起点坐标

		//如果当前落子角色不是黑子，说明机器落子，人不能点击
		if obj.currentRole != Black {
			return
		}

		// 要保证点击点在棋盘范围里面
		if obj.x >= obj.startX && obj.x <= obj.startX+8*obj.gridW && obj.y >= obj.startY && obj.y <= obj.startX+8*obj.gridH {
			// 棋盘的位置转换转换为坐标下标值
			i := (obj.x - obj.startX) / obj.gridW
			j := (obj.y - obj.startY) / obj.gridH

			//保证i, j在0~7范围里
			if i >= 0 && i <= 7 && j >= 0 && j <= 7 {
				//fmt.Printf("i = %d, j = %d\n", i, j)

				if obj.JudgeRule(i, j, obj.currentRole, true) > 0 { //能吃子才更新棋盘
					obj.window.QueueDraw() //刷新绘图区域
					obj.ChangeRole()       //改变落子角色
				}
			}
		}

	}
}

//方法，事件、信号处理，回调函数如果简单使用匿名函数，否则自定义函数
func (obj *Chessboard) HandleSignal() {

	//鼠标按下事件处理，MousePressEvent为其回调函数，把obj传递给回调函数
	obj.window.Connect("button-press-event", MousePressEvent, obj)

	//鼠标移动事件处理，实现窗口的移动
	obj.window.Connect("motion-notify-event", func(ctx *glib.CallbackContext) {
		//获取鼠标属性结构体变量，系统内部的变量，不是用户传参变量
		arg := ctx.Args(0)
		//还是EventButton
		event := *(**gdk.EventButton)(unsafe.Pointer(&arg))

		//移动窗口
		obj.window.Move(int(event.XRoot)-obj.x, int(event.YRoot)-obj.y)
	})

	//改变窗口大小时，触发"configure-event"，然后手动刷新绘图区域，否则图片会重叠
	obj.window.Connect("configure-event", func() {
		obj.window.QueueDraw() //刷新绘图区域
	})

	//绘图（曝光）事件，其回调函数PaintEvent做绘图操作，把obj传递给回调函数
	obj.window.Connect("expose-event", PaintEvent, obj)

	//最小化按钮信号处理
	obj.buttonMin.Connect("clicked", func() {
		obj.window.Iconify() //窗口最小化
	})

	//关闭按钮信号处理
	obj.buttonClose.Connect("clicked", func() {
		gtk.MainQuit() //程序结束
	})
}

//定时器处理函数，角色提示，达到一闪一闪的效果
func ShowTip(obj *Chessboard) {
	if obj.currentRole == Black { //当前该黑子下
		//白子提示图片隐藏
		obj.imageWhite.Hide()

		if obj.imageBlack.GetVisible() == true { //原来显示的，则隐藏
			obj.imageBlack.Hide()
		} else {
			obj.imageBlack.Show()
		}

	} else if obj.currentRole == White { //当前该白子下
		//黑子提示图片隐藏
		obj.imageBlack.Hide()

		if obj.imageWhite.GetVisible() == true { //原来显示的，则隐藏
			obj.imageWhite.Hide()
		} else {
			obj.imageWhite.Show()
		}
	}
}

//方法，主要做些初始化操作
func (obj *Chessboard) InitChess() {
	//当前落子角色
	obj.currentRole = Black
	//obj.currentRole = White

	//全部标记为Empty
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			obj.chess[i][j] = Empty
		}
	}

	//中间位置
	obj.chess[3][3] = Black
	obj.chess[4][4] = Black
	obj.chess[4][3] = White
	obj.chess[3][4] = White

	//刷新绘图区域
	obj.window.QueueDraw()

	//先隐藏提示图片
	obj.imageBlack.Hide()
	obj.imageWhite.Hide()

	//黑白子个数各自为2
	obj.labelBlack.SetText("2")
	obj.labelWhite.SetText("2")

	//落子提示
	obj.tipTimerId = glib.TimeoutAdd(500, func() bool {
		ShowTip(obj) //调用函数处理
		return true
	})

	//倒计时
	obj.timeNum = 20
	obj.labelTime.SetText(strconv.Itoa(obj.timeNum)) //标签显示时间

	//启动倒计时定时器
	obj.leftTimerId = glib.TimeoutAdd(1000, func() bool {
		obj.timeNum--
		obj.labelTime.SetText(strconv.Itoa(obj.timeNum)) //标签显示时间
		if obj.timeNum == 0 {                            //时间到
			obj.ChangeRole() //改变落子角色
		}

		return true
	})
}

func main() {
	gtk.Init(&os.Args)

	var obj Chessboard //创建结构体变量

	//创建控件，设计属性
	window := obj.CreateWindow()

	//初始化数据
	obj.InitChess()

	//事件、信号处理
	obj.HandleSignal()

	//显示控件
	window.Show()

	gtk.Main()
}
