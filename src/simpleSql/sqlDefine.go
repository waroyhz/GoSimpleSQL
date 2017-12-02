package sSql

var ALL string = " * "
var MAX string = "max(%s)"
var COUNT string = " count(*) "
var PARAM string = "?"
var InPARAM string = " in(?) "
var NotInPARAM string = " not in(?) "
var LEFT string = "("
var RIGHT string = ") "
var COMMA string = ","

const (
	EQ string = "="
	NE string = "!="
	GT string = ">"
	GE string = ">="
	LT string = "<"
	LE string = "<="
	BETWEEN string = " between "
	AND string = " and "
	OR string = " or "

	ISNOTNULL string = " IS NOT NULL "
	ISNULL 	string = " IS NULL "
	ORDERBY string = " order by "
	DESC string = " DESC "
	ASC string = " ASC "
	LIMIT string = " LIMIT "
	LIKE  string = " LIKE "
	LEFTJOIN string = " left join "
	UNIX_TIMESTAMP = "unix_timestamp"


	selecT string = "select "
	update string = "update "
	delete string = "delete from "
	insert string = "insert into "
	set string = " set "
	from string = " from "
	where string = " where "
	group string = " group by "
	order string = " order by "
	on string = " on "
	values string = " values"

	inparam string = "?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?"
	tablep = "abcde"
)

type SQL struct {
	intParamAllowZero bool //可以添加值为0的int类型
	query             string
	from              string
	where             string
	group             string
	order             string
	on                string
	args              []interface{}

	isselect          bool
	isupdate          bool
	isdelete          bool
	isinsert          bool
}

type Set struct {
	Field string
	Value string
}

type WhereEQ struct {
	Field string
	Value string
}
