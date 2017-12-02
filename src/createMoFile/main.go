package main

import (
	"bufio"
	"fmt"
	"log"

	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"tools/shell"
)

type modules struct {
	//这里引入model
	//user   mhj_models.UserEntity

}

var pathMO = "/src/mhj/models/modelsDefine"
var path string

func main() {
	fmt.Println("模型枚举生成工具!")
	//fmt.Printf("当前路径：%s\n", getCurrentDirectory())

	//fmt.Println(command("pwd"))

	path = shell.GetCurrentPath()
	spath:= strings.Split(path,"/")
	spath[len(spath)-1]="smartShared"
	path=strings.Join(spath,"/")
	fmt.Printf("当前路径：%s\n", path)
	path = strings.Replace(path, "\n", "", -1) + pathMO

	fmt.Printf("检测路径：%s\n", path)
	if !shell.IsDirExists(path) {
		fmt.Printf("未找到modelsDefine文件夹，请确认路径是否正确！\n")
		return
	}

	m := modules{}
	s := reflect.ValueOf(&m).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		createMOFile(typeOfT.Field(i).Type.String(), f)
	}
}

var digitsRegexp = regexp.MustCompile(`(\d+)\D+(\d+)`)

func createMOFile(typename string, m reflect.Value) {
	typeOfT := m.Type()
	classname := substr(typename, strings.Index(typename, ".") + 1, 0xffff)

	filename := path + "/mo" + classname + ".go"
	if shell.IsFileExists(filename) {
		os.Remove(filename)
	}
	os.Remove(filename)
	fmt.Printf("写入文件：%s\n", filename)
	if f, err1 := os.OpenFile(filename, os.O_CREATE | os.O_RDWR, 0666); err1 != nil {
		fmt.Printf("创建文件%s失败\n", filename)
	} else {

		w := bufio.NewWriter(f) //创建新的 Writer 对象

		if _, err3 := w.WriteString("package modelsDefine\n"); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}

		w.WriteString("/*\n		模型枚举文件，请勿修改！\n*/\n")

		//linestring := ""
		fieldName := ""
		for i := 0; i < m.NumField(); i++ {

			fmt.Print("TypeName=", typeOfT.Field(i).Type.Name())
			fieldName=typeOfT.Field(i).Tag.Get("field")
			//tag:=typeOfT.Field(i).Tag
			//tags := strings.Split(string(tag), "\"")
			//if len(tags) > 1 {
			//	fieldName = tags[1]
			//
			//} else {
			//	fieldName = typeOfT.Field(i).Name
			//
			//}
			writeField(classname, typeOfT.Field(i), "", fieldName, w, filename)

		}
		w.Flush()
		f.Close()
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func createChildField(classname, parentFieldName string, childtype reflect.Type, w *bufio.Writer, filename string) {
	s := reflect.New(childtype).Elem()
	typeOfT := childtype
	fmt.Println(s, typeOfT)
	fieldName := ""
	for i := 0; i < s.NumField(); i++ {

		tags := strings.Split(string(typeOfT.Field(i).Tag), "\"")
		if len(tags) > 1 {
			fieldName = tags[1]

		} else {
			fieldName = typeOfT.Field(i).Name

		}

		writeField(classname, typeOfT.Field(i), "", parentFieldName + fieldName, w, filename)

	}
}

func writeField(classname string, field reflect.StructField, parentFieldName, fieldName string, w *bufio.Writer, filename string) {
	linestring := ""
	if field.Type.String() == "mgo.DBRef" {
		//联级查询
		linestring = fmt.Sprintf("const Mo%s_%s string = \"%s\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}

		linestring = fmt.Sprintf("const Mo%s_%s_ref string = \"%s.$ref\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}
		fmt.Printf("写入内容： %s", linestring)
		linestring = fmt.Sprintf("const Mo%s_%s_id string = \"%s.$id\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}
		fmt.Printf("写入内容： %s", linestring)
		linestring = fmt.Sprintf("const Mo%s_%s_db string = \"%s.$db\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}
		fmt.Printf("写入内容： %s", linestring)
	} else {
		linestring = fmt.Sprintf("const Mo%s_%s string = \"%s\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}
		fmt.Printf("写入内容： %s", linestring)

		typename := fmt.Sprint(field.Type)
		if strings.Contains(typename, "Tile") {
			fmt.Println(typename)
		}
		if strings.Contains(typename, "mhj_models") {
			if parentFieldName == "" {
				parentFieldName = fieldName + "."
			} else {
				parentFieldName = parentFieldName + fieldName + "."
			}

			if strings.Index(typename, "*") == 0 {
				createChildField(classname + "_" + field.Name, parentFieldName, field.Type.Elem(), w, filename)
			} else {
				createChildField(classname + "_" + field.Name, parentFieldName, field.Type, w, filename)
			}
		}
	}
}
