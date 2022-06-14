package resource

import (
	"github.com/infraboard/mcube/http/request"
	_ "github.com/infraboard/mcube/http/request"
	"net/http"
	"strings"
)

const (
	AppName = "resource"
)

//定义四种表达式，用于匹配values的场景2 使用符号搜索
//tag的比较操作符号类比promethus做的	官网找到四种操作符号
//把四种符号映射为sql语句
type Operator string

const (
	// SQL 比较操作  =
	Operator_EQUAL = "="
	// SQL 比较操作  !=
	Operator_NOT_EQUAL = "!="
	// SQL 比较操作  LIKE
	Operator_LIKE_EQUAL = "=~"
	// SQL 比较操作  NOT LIKE
	Operator_NOT_LIKE_EQUAL = "!~"
)

// HasTag 扩展SearchRequest 给他一个方法HasTag，让他返回一个当前的查询参数到底有没有带着tag来做一个查询
//判断下长度是不是大于0，如果大于0就可以做过滤，如果等于0就没有tag
func (r *SearchRequest) HasTag() bool {
	return len(r.Tags) > 0
}

//围绕Operator 需要扩展Tagselectot
//TagSelector 映射为relastionship
//对多个值的比较的关系做一个说明
//比如说你传的是一个叫做app=～app1,app2  是一种白名单策略（包含策略）
//表示这个app的值为tag，key为app 。这一个key的值等于app1或者app2，都是可以的
//你不能说app1和app2 是and 的关系，一定是一个or的关系
//app!=app3,app4,tag_key=app tag_vlue NOT LIKE (app3,app4),是一种黑名单策略（排除策略）

func (s *TagSelector) Relationship() string {
	switch s.Operator {
	//如果是一个like的等于符号就给予一个or的条件。
	case Operator_EQUAL, Operator_LIKE_EQUAL:
		return " OR "
	case Operator_NOT_EQUAL, Operator_NOT_LIKE_EQUAL:
		return "AND"
	default:
		return "OR"
	}

}

func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		Items: []*Resource{},
	}
}

//add逻辑就是
func (s *ResourceSet) Add(item *Resource) {
	//加入到items列表里面
	s.Items = append(s.Items, item)
}

func NewDefaultResource() *Resource {
	return &Resource{
		Base:        &Base{},
		Information: &Information{},
	}
}

//从数据库取得数据需要进行格式化
//独立绑定一个方法单独处理这个逻辑
func (i *Information) LoadPrivateIPString(s string) {
	if s != "" {
		//按照逗号分割出PrivateIp，PrivateIp是一个ip地址的列表
		i.PrivateIp = strings.Split(s, ",")
	}
}

func (i *Information) LoadPublicIPString(s string) {
	if s != "" {
		//PublicIp 也是一个ip地址的列表
		i.PublicIp = strings.Split(s, ",")
	}
}

//NewDefaultTag 方法
//默认给予的tag 是TagType_USER 用户自定义标签, 允许用户修改
func NewDefaultTag() *Tag {
	return &Tag{
		Type:   TagType_USER,
		Weight: 1,
	}
}

// ResourceIds 方法
//把item里面的所有id取出，放到ids []string 数组里面去返回
// 把属性里面的id取出来转换为string的一个操作
//循环到append到一个数组里面
func (s *ResourceSet) ResourceIds() (ids []string) {
	for i := range s.Items {
		ids = append(ids, s.Items[i].Base.Id)
	}

	return
}

//UpdateTag 是一个循环
//先循环所有的tag 把tag里面的 ResourceId 跟我们的列表里面的item里面的id相等
//我们就把这个tag添加到information里面Tag 字段里面
//所以我们的information 需要添加一个addtag的功能
func (s *ResourceSet) UpdateTag(tags []*Tag) {
	for i := range tags {
		for j := range s.Items {
			if s.Items[j].Base.Id == tags[i].ResourceId {
				s.Items[j].Information.AddTag(tags[i])
			}
		}
	}
}

// 把tag添加到information的哪个字段上面
func (r *Information) AddTag(t *Tag) {
	r.Tags = append(r.Tags, t)
}

// keywords=xx&domain=xx&tag=app=~app1,app2,app3
func NewSearchRequestFromHTTP(r *http.Request) (*SearchRequest, error) {
	qs := r.URL.Query()
	req := &SearchRequest{
		Page:        request.NewPageRequestFromHTTP(r),
		Keywords:    qs.Get("keywords"),
		ExactMatch:  qs.Get("exact_match") == "true",
		Domain:      qs.Get("domain"),
		Namespace:   qs.Get("namespace"),
		Env:         qs.Get("env"),
		Status:      qs.Get("status"),
		SyncAccount: qs.Get("sync_account"),
		WithTags:    qs.Get("with_tags") == "true",
		Tags:        []*TagSelector{},
	}

	umStr := qs.Get("usage_mode")
	if umStr != "" {
		mode, err := ParseUsageModeFromString(umStr)
		if err != nil {
			return nil, err
		}
		req.UsageMode = &mode
	}

	rtStr := qs.Get("resource_type")
	if rtStr != "" {
		rt, err := ParseTypeFromString(rtStr)
		if err != nil {
			return nil, err
		}
		req.Type = &rt
	}

	// 单独处理Tag参数 app~=app1,app2,app3 --> TagSelector ---> req
	tgStr := qs.Get("tag")
	if tgStr != "" {
		tg, err := NewTagsFromString(tgStr)
		if err != nil {
			return nil, err
		}
		req.AddTag(tg...)
	}

	return req, nil
}
