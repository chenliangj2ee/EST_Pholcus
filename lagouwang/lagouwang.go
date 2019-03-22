package pholcus_lib

// 基础包
import (
	// "github.com/henrylee2cn/pholcus/common/goquery" //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	. "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	"github.com/henrylee2cn/pholcus/logs"                   //信息输出

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	"encoding/json"

	// 字符串处理包
	"regexp"
	"strings"
	// 其他包
	"fmt"
	// "math"
	// "time"
)

func init() {
	LaGou.Register()
}

var LaGou = &Spider{
	Name:        "陈亮拉勾网",
	Description: "陈亮拉勾网 [https://www.lagou.com]",
	// Pausetime: 300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "https://www.lagou.com/jobs/list_?labelWords=&fromSearch=true&suginput=", Rule: "获取版块URL"})
		},

		Trunk: map[string]*Rule{

			"获取版块URL": {
				ParseFunc: func(ctx *Context) {

					pageStr := ctx.GetText()
					cityReg := regexp.MustCompile(`<a rel="nofollow" href="javascript:;" data-lg-tj-id="8p00" data-lg-tj-no="
                                                            ([.]*)
                                                    " data-lg-tj-cid="idnull" data-id="([.]*)" class="more-city-name">([.]*)</a>`)

					result := cityReg.FindAllStringSubmatch(pageStr, -1)
					fmt.Println(result)

				},
			},

			"搜索结果": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					src := query.Find("script").Text()

					re, _ := regexp.Compile(`"auctions".*,"recommendAuctions"`)
					src = re.FindString(src)

					re, _ = regexp.Compile(`"auctions":`)
					src = re.ReplaceAllString(src, "")

					re, _ = regexp.Compile(`,"recommendAuctions"`)
					src = re.ReplaceAllString(src, "")

					re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
					// src = re.ReplaceAllStringFunc(src, strings.ToLower)
					src = re.ReplaceAllString(src, " ")

					src = strings.Trim(src, " \t\n")

					infos := []map[string]interface{}{}

					err := json.Unmarshal([]byte(src), &infos)

					if err != nil {
						logs.Log.Error("error is %v\n", err)
						return
					} else {
						for _, info := range infos {
							ctx.AddQueue(&request.Request{
								Url:  "http:" + info["detail_url"].(string),
								Rule: "商品详情",
								Temp: ctx.CreatItem(map[int]interface{}{
									0: info["raw_title"],
									1: info["view_price"],
									2: info["view_sales"],
									3: info["nick"],
									4: info["item_loc"],
								}, "商品详情"),
								Priority: 1,
							})
						}
					}
				},
			},
			"商品详情": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"标题",
					"价格",
					"销量",
					"店铺",
					"发货地",
				},
				ParseFunc: func(ctx *Context) {
					r := ctx.CopyTemps()

					re := regexp.MustCompile(`"newProGroup":.*,"progressiveSupport"`)
					d := re.FindString(ctx.GetText())

					if d == "" {
						h, _ := ctx.GetDom().Find(".attributes-list").Html()
						d = UnicodeToUTF8(h)
						d = strings.Replace(d, "&nbsp;", " ", -1)
						d = CleanHtml(d, 5)
						d = strings.Replace(d, "产品参数：\n", "", -1)

						for _, v := range strings.Split(d, "\n") {
							if v == "" {
								continue
							}
							feild := strings.Split(v, ":")
							// 去除英文空格
							// feild[0] = strings.Trim(feild[0], " ")
							// feild[1] = strings.Trim(feild[1], " ")
							// 去除中文空格
							feild[0] = strings.Trim(feild[0], " ")
							feild[1] = strings.Trim(feild[1], " ")

							if feild[0] == "" || feild[1] == "" {
								continue
							}

							ctx.UpsertItemField(feild[0])
							r[feild[0]] = feild[1]
						}

					} else {
						d = strings.Replace(d, `"newProGroup":`, "", -1)
						d = strings.Replace(d, `,"progressiveSupport"`, "", -1)

						infos := []map[string]interface{}{}

						err := json.Unmarshal([]byte(d), &infos)

						if err != nil {
							logs.Log.Error("error is %v\n", err)
							return
						} else {
							for _, info := range infos {
								for _, attr := range info["attrs"].([]interface{}) {
									a := attr.(map[string]interface{})
									ctx.UpsertItemField(a["name"].(string))
									r[a["name"].(string)] = a["value"]
								}
							}
						}
					}

					ctx.Output(r)
				},
			},
		},
	},
}
