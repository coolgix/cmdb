package resource

import (
	_ "github.com/infraboard/mcube/http/request"
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

func NewDefaultTag() *Tag {
	return &Tag{
		Type:   TagType_USER,
		Weight: 1,
	}
}

func (s *ResourceSet) ResourceIds() (ids []string) {
	for i := range s.Items {
		ids = append(ids, s.Items[i].Base.Id)
	}

	return
}

func (r *Information) AddTag(t *Tag) {
	r.Tags = append(r.Tags, t)
}

func (s *ResourceSet) UpdateTag(tags []*Tag) {
	for i := range tags {
		for j := range s.Items {
			if s.Items[j].Base.Id == tags[i].ResourceId {
				s.Items[j].Information.AddTag(tags[i])
			}
		}
	}
}
