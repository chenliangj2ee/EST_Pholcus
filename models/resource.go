package models

type Resource struct {
	Title    string  `json:"title"`     //标题
	Image    string  `json:"image"`     //图片
	Des      string  `json:"des"`       //描述
	Link     string  `json:"link"`      //点击链接
	ClickNum int     `json:"click_num"` //点击量
	Price    float64 `json:"price"`     //价格
	Author   string  `json:"author"`    //作者
	ClassNum int     `json:"class_num"` //课时
	Remark   string  `json:"remark"`    //备注
}

func (res *Resource) InsertDB() {

}
