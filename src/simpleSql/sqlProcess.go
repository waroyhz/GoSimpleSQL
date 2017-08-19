package sSql

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

func isEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
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