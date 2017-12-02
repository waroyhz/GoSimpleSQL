package sql

import (
	"strings"
	"reflect"
)

func query(args ... string) string {
	sql := ""
	for _, str := range args {
		sql += str
	}
	return sql
}


func querySelect(args ... string) string {
	sql := []string{}
	for _, str := range args {
		sql = append(sql, str)
	}
	return strings.Join(sql, ",")
}

func isEmpty(name string,tag reflect.StructTag,a int64) bool {
	//v := reflect.ValueOf(a)
	//if v.Kind() == reflect.Ptr {
	//	v = v.Elem()
	//}

	//defvalue:=reflect.Zero(v.Type()).Interface()
	//nv:=v.Interface()
	ret:= a==0

	nullYES:=tag.Get("null")=="YES"
	keyMUL:=tag.Get("key")=="MUL"

	if ret && !nullYES {
		if 0 ==a{
			if name=="Id"{ //id为0时置为空
				ret=true
			}else{
				ret= false
			}
		}
	}else if keyMUL && nullYES {
		if 0==a{
			ret=true
		}
	}
	//if ret{
	//	log4go.Info(name,a,v,reflect.Zero(v.Type()).Interface(),reflect.Zero(v.Type()),ret)
	//}
	return ret
}

func getInsertQuery(cols string) string {
	var params []string
	i := len(strings.Split(cols, ","))
	for j := 0; j < i; j++ {
		params = append(params, PARAM)
	}
	param := strings.Join(params, ",")
	return " " + LEFT + cols + RIGHT + values + LEFT + param + RIGHT
}