package csdn

import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"regexp"
	"pholcus/utils"
	"strconv"
	"github.com/henrylee2cn/pholcus/common/goquery"
	"pholcus/models"
)

var (
	host      = "https://www.imooc.com"
	url01     = "https://www.imooc.com/course/list"
	pageTotal int       //总页数，默认29
	lenght    int       //总数量
	finish    chan bool //记录是否爬虫结束
)

func init() {
	CSDN.Register()
	utils.Mylog("CSDN学院注册...")
}

var CSDN = &Spider{
	Name:         "CSDN学院",
	Description:  "CSDN学院 [https://www.imooc.com]",
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			lenght = 0
			utils.Mylog("CSDN学院-开始爬虫.....")
			ctx.AddQueue(&request.Request{Url: url01, Rule: "解析总页数"})

		},

		Trunk: map[string]*Rule{
			"解析总页数":  {ParseFunc: ParsePageTotal,},
			"解析课程列表": {ParseFunc: ParseList,},
		},
	},
}

/*
1、解析当前页，获取最大页数
2、循环页码，分页加载数据
*/
var ParsePageTotal = func(ctx *Context) {
	pageStr := ctx.GetText()
	cityReg := regexp.MustCompile(`page=([1-9]*)">尾页</a>`)
	result := cityReg.FindAllStringSubmatch(pageStr, -1)
	page, err := strconv.Atoi(result[0][1])
	if err != nil {
		pageTotal = 29
	}
	pageTotal = page

	for i := 1; i <= pageTotal; i++ {
		url := url01 + "?page=" + strconv.Itoa(i)
		ctx.AddQueue(&request.Request{Url: url, Rule: "解析课程列表"})
	}
	finishFun(pageTotal)

}

/*
1、解析当前页所有item
2、保存到数据库
*/
var ParseList = func(ctx *Context) {
	query := ctx.GetDom()
	resultSel := query.Find(".course-card-container")
	//判断是否有子节点
	if len(resultSel.Nodes) > 0 {
		//变量循环每个子节点
		resultSel.Each(func(i int, selection *goquery.Selection) {

			link, _ := selection.Find(".course-card").Attr("hre2f")  //点击链接
			image, _ := selection.Find(".course-banner").Attr("src") //图片地址
			title := selection.Find(".course-card-name").Text()      //课程标题
			des := selection.Find(".course-card-desc").Text()        //课程描述

			if "" == link || "" == title || "" == image {
				utils.Mylog("CSDN学院-解析数据有误，请检查.....")
				utils.Mylog("link:", link, "title:", title, "image:", image)
			} else {
				lenght++
				res := models.Resource{Title: title, Image: image, Link: host + link, Des: des, ClickNum: 0}
				res.InsertDB()
			}

			if i == len(resultSel.Nodes)-1 {
				finish <- true
			}

		})
	} else {
		finish <- true
		utils.Mylog("CSDN学院-解析数据失败.....")
	}

}

func finishFun(pageNum int) {
	finish = make(chan bool, pageNum)
	for i := 0; i < pageNum; i++ {
		<-finish
	}
	utils.Mylog("CSDN学院-爬虫结束.....", lenght)
}
