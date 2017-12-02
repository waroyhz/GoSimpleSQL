package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"time"
	"tools/shell"
)

/*
数据库表生成模型的工具
*/


var tableNames = []string{
	//数据库的表名
	"user_entity",

}

const (
	driverName = "mysql"
	database = "mhj:mhj123@tcp(192.168.5.15:3306)/mhj?charset=utf8"
	childpathModels = "/src/mhj/models"
	childpathMO = "/src/mhj/models/modelsDefine"
)
var pathMO string
var pathModels string

type Colum struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

func main() {
	fmt.Println("模型枚举生成工具!")
	//fmt.Printf("当前路径：%s\n", getCurrentDirectory())

	//fmt.Println(command("pwd"))

	tpathMO:= shell.GetCurrentPath()

	spath:= strings.Split(tpathMO,"/")
	spath[len(spath)-1]="smartShared"
	tpathMO=strings.Join(spath,"/")

	fmt.Printf("当前路径：%s\n", pathMO)
	pathMO = strings.Replace(tpathMO, "\n", "", -1) + childpathMO
	pathModels = strings.Replace(tpathMO, "\n", "", -1) + childpathModels

	fmt.Printf("Models路径：%s\n", pathModels)
	fmt.Printf("检测路径：%s\n", pathMO)
	if !shell.IsDirExists(pathMO) {
		fmt.Printf("未找到modelsDefine文件夹，请确认路径是否正确！\n")
		return
	}

	if db, err := sql.Open(driverName, database); err == nil {
		for _, tableName := range tableNames {
			rows, rowerr := db.Query("show columns from " + tableName)
			if rowerr == nil {
				colums := []Colum{}
				for rows.Next() {
					rowData := Colum{}
					//rows.Columns()
					rows.Scan(&rowData.Field, &rowData.Type, &rowData.Null, &rowData.Key, &rowData.Default, &rowData.Extra)
					//if serr!=nil {
					//	log.Fatal(err)
					//}
					fmt.Println(rowData)
					colums = append(colums, rowData)
				}
				createModuleFile(pathModels,tableName, colums)
				//createMOFile(pathMO,tableName,colums)
			} else {
				log.Fatal(rowerr)
			}
		}

	} else {
		log.Fatal(err)
	}

	//m := modules{}
	//s := reflect.ValueOf(&m).Elem()
	//typeOfT := s.Type()
	//for i := 0; i < s.NumField(); i++ {
	//	f := s.Field(i)
	//	createMOFile(typeOfT.Field(i).Type.String(), f)
	//}
}
func createModuleFile(modulePath string,tableName string, colums []Colum) {
	modelName := formatName(tableName)

	sourceFile := `package mhj_models
/**
自动生成模型工具
			by waroy
*/
import (
	%%import%%
)

const TABLE_%%TABLE%% string = "%%table%%"

type %%TABLE%% struct {
%%FIELDS%%
}
	`

	typeimport :=map[string]string{}
	fields := []string{}
	for _, c := range colums {
		field := fmt.Sprintf("	%s	%s	%s	`%s`", formatName(c.Field), formatType(c.Type, typeimport), "",fmt.Sprintf(`field:"%s" key:"%s" type:"%s" null:"%s" default:"%s" extra:"%s"`,c.Field,c.Key,c.Type,c.Null,c.Default,c.Extra))
		fields = append(fields, field)
	}
	strimports:=[]string{}
	for _,strimport:=range typeimport{
		strimports=append(strimports,fmt.Sprintf("\"%s\"",strimport))
	}

	strimport:=strings.Join(strimports,"\n")


	fieldData := strings.Join(fields, "\n")

	filedata := strings.Replace(sourceFile, "%%TABLE%%", modelName, -1)
	filedata = strings.Replace(filedata, "%%table%%", tableName, -1)
	filedata = strings.Replace(filedata, "%%FIELDS%%", fieldData, -1)
	filedata = strings.Replace(filedata, "%%import%%", strimport, -1)

	//println(modelName, filedata)
	modulefile:=modulePath+"/"+modelName+".go"
	fmt.Printf("写入文件：%s\n",modulefile )
	os.Remove(modulefile)
	if f, err1 := os.OpenFile(modulefile, os.O_CREATE | os.O_RDWR, 0666); err1 != nil {
		fmt.Printf("创建文件%s失败\n", modulefile)
	} else {

		w := bufio.NewWriter(f) //创建新的 Writer 对象
		w.WriteString(filedata)
		w.Flush()
		f.Close()
	}
}
func formatType(sqlType string, typeimport map[string]string) string {
	strType := sqlType
	li := strings.Index(sqlType, "(")
	if li > 0 {
		strType = string([]byte(sqlType)[:li])
	}
	switch strType {
	case "int":
		return reflect.TypeOf(int(1)).Name()
	case "longtext", "varchar","text":
		return reflect.TypeOf(string("")).Name()
	case "char":
		//return "[]byte"//+reflect.TypeOf(byte('1')).Name()
		return reflect.TypeOf(string("")).Name()
	case "bigint":
		return reflect.TypeOf(int64(1)).Name()
	case "tinyint":
		return reflect.TypeOf(bool(true)).Name()
	case "datetime":
		typeimport[strType]=reflect.TypeOf(time.Now()).PkgPath()
		return reflect.TypeOf(time.Now()).PkgPath()+"."+reflect.TypeOf(time.Now()).Name()
	default:
		log.Fatalf("类型 %s 没有 转换！！！",strType)
	}
	return ""
}
func formatName(name string) string {
	lowname:=strings.ToLower(name);
	fmtName := lowname;
	//if strings.Contains(fmtName, "_") {
	strs := strings.Split(lowname, "_")
	for i, s := range strs {
		ups := strings.ToUpper(s)
		strs[i] = string([]byte(ups)[0:1]) + string([]byte(s)[1:])
	}
	fmtName = strings.Join(strs, "")
	//}
	return fmtName
}

