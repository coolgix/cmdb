package resource

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
func (s *TagSelector) Relationship() string {
	switch s.Operator {
	//如果是一个like的等于符号就给予一个or的条件
	case Operator_EQUAL, Operator_LIKE_EQUAL:
		return " OR "
	case Operator_NOT_EQUAL, Operator_NOT_LIKE_EQUAL:
		return "AND"
	default:
		return "OR"
	}

}
