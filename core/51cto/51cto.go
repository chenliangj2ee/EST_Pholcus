package _1cto

import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"regexp"
	"pholcus/utils"
	"strconv"
	"github.com/henrylee2cn/pholcus/common/goquery"
)

var (
	host      = "https://edu.csdn.net"
	url01     = "https://edu.csdn.net/courses/k"
	pageTotal int       //总页数，默认339
	lenght    int       //总数量
	finish    chan bool //记录是否爬虫结束
)

func init() {
	CTO51.Register()
	utils.Mylog("51CTO学院注册...")
}

var CTO51 = &Spider{
	Name:         "51CTO学院",
	Description:  "51CTO学院 [https://edu.csdn.net]",
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			lenght = 0
			utils.Mylog("51CTO学院-开始爬虫.....")
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
	cityReg := regexp.MustCompile(`<a class="btn btn-xs btn-default" href="https://edu.csdn.net/courses/k/p([1-9]*)">([1-9]*)</a>`)
	result := cityReg.FindAllStringSubmatch(pageStr, -1)
	page, err := strconv.Atoi(result[0][1])
	if err != nil {
		pageTotal = 29
	}
	pageTotal = page

	for i := 1; i <= pageTotal; i++ {
		url := url01 + "" + strconv.Itoa(i)
		ctx.AddQueue(&request.Request{Url: url, Rule: "解析课程列表"})
	}
	finishFun(pageTotal)

}

//<div class="course_item">
//<a href="https://edu.csdn.net/course/detail/16868" target="_blank">
//<dl class="course_item_inner clearfix">
//<dd>
//<span class="course_lessons">87课时</span>
//<img src="https://img-bss.csdn.net/201913017329703_30493.jpg" width="179" height="120" alt="深度学习与PyTorch入门实战教程">
//</dd>
//<dt>
//<div class="titleInfor">
//<span class="title ellipsis-2">
//深度学习与PyTorch入门实战教程                        </span>
//</div>
//<p class="subinfo">
//<span class="num">
//<em>￥4.57</em>/课时                        </span>
//<!--                        <span class="num">--><!--课时</span>-->
//<span class="lecname ellipsis" title="龙良曲">龙良曲</span></p>
//<p class="priceinfo clearfix">
//<i>
//￥398.00                        </i>
//</p>
//</dt>
//</dl>
//</a>
//</div>
/*
1、解析当前页所有item
2、保存到数据库
*/
var ParseList = func(ctx *Context) {
	query := ctx.GetDom()
	resultSel := query.Find(".course_item")
	//判断是否有子节点
	if len(resultSel.Nodes) > 0 {
		//变量循环每个子节点
		resultSel.Each(func(i int, selection *goquery.Selection) {
			link, _ := selection.Find("a").Attr("href")
			image, _ := selection.Find("img").Attr("src")
			classNum := selection.Find(".course_lessons").Text() //课时
			title := selection.Find(".title").Text()             //课时
			price := selection.Find(".priceinfo").Text()         //课程标题
			author := selection.Find(".lecname").Text()          //课程标题

			classNum = utils.Trim(classNum)
			title = utils.Trim(title)
			price = utils.Trim(price)
			author = utils.Trim(author)

			utils.Mylog(link, image, classNum, title, price, author)

			//if "" == link || "" == title || "" == image {
			//	utils.Mylog("51CTO学院-解析数据有误，请检查.....")
			//	utils.Mylog("link:", link, "title:", title, "image:", image)
			//} else {
			//	lenght++
			//	res := models.Resource{Title: title, Image: image, Link: host + link, Des: des, ClickNum: 0}
			//	res.InsertDB()
			//}
			//
			//if i == len(resultSel.Nodes)-1 {
			//	finish <- true
			//}

		})
	} else {
		finish <- true
		utils.Mylog("51CTO学院-解析数据失败.....")
	}

}

func finishFun(pageNum int) {
	finish = make(chan bool, pageNum)
	for i := 0; i < pageNum; i++ {
		<-finish
	}
	utils.Mylog("51CTO学院-爬虫结束.....", lenght)
}
