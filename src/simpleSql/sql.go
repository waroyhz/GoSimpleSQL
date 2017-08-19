package sSql

import (
	"fmt"
	"strings"
	"log"
	"reflect"
	"database/sql"
	"time"
	"github.com/alecthomas/log4go"
	"unicode"
)

func NewCommand(tablename string) *SQL {
	s := new(SQL)
	s.from = tablename
	return s
}

func (this*SQL)Select(args ... string) *SQL {
	this.query = querySelect(args...)
	this.isselect = true
	return this
}

func (this*SQL)Delete(args ... WhereEQ) *SQL {
	strwhere := []string{}
	for _, set := range args {
		strwhere = append(strwhere, fmt.Sprintf("%s=%s", FieldFormat(set.Field), set.Value))
	}
	this.where = where + strings.Join(strwhere, AND)
	this.isdelete = true
	return this
}

//用于同一个字段
func (this *SQL)Deletes() *SQL {
	this.isdelete = true
	return this
}

func (this*SQL)Update(args ... Set) *SQL {
	strsets := []string{}
	for _, set := range args {
		strsets = append(strsets, fmt.Sprintf("%s=%s", set.Field, set.Value))
	}
	this.query = strings.Join(strsets, ",")
	this.isupdate = true
	return this
}

func (this*SQL)IntParamAllowZero() *SQL {
	this.intParamAllowZero = true
	return this
}

func (this*SQL)Insert(obj interface{}) *SQL {
	strsets := []string{}
	objType := reflect.TypeOf(obj).Elem()
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objType.NumField(); i++ {
		var value interface{}
		//获取field的name
		fieldType := objType.Field(i).Type.Kind()   //field的类型
		fieldName := objType.Field(i).Name             //field的name
		fieldValue := objValue.FieldByName(fieldName)

		switch fieldType {
		case reflect.String:
			if ns, ok := fieldValue.Interface().(sql.NullString); ok {
				value = nil
				if ns.Valid {
					value = ns.String
				}
			} else if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					value = nil
				} else {
					value = fieldValue.Elem().String()
				}
			} else {
				value = fieldValue.String()
			}
			if value != nil&&value != "" {
				strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
				this.args = append(this.args, value)
			}

		case reflect.Bool:
			if nb, ok := fieldValue.Interface().(sql.NullBool); ok {
				value = nil
				if nb.Valid {
					value = nb.Bool
				}
			} else if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					value = nil
				} else {
					value = fieldValue.Elem().Bool()
				}
			} else {
				value = fieldValue.Bool()
			}
			if value != nil {
				strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
				this.args = append(this.args, value)
			}

		case reflect.Int, reflect.Int32, reflect.Int64:
			if ni, ok := fieldValue.Interface().(sql.NullInt64); ok {
				value = nil
				if ni.Valid {
					value = ni.Int64
				}
			} else if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					value = nil
				} else {
					value = fieldValue.Elem().Int()
				}
			} else {
				value = fieldValue.Int()
			}
			if this.intParamAllowZero {
				if value != nil {
					strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
					this.args = append(this.args, value)
				}
			} else {
				if value != nil && !isEmpty(value) {
					strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
					this.args = append(this.args, value)
				}
			}

		case reflect.Float32, reflect.Float64:
			if ni, ok := fieldValue.Interface().(sql.NullFloat64); ok {
				value = nil
				if ni.Valid {
					value = ni.Float64
				}
			} else if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					value = nil
				} else {
					value = fieldValue.Elem().Int()
				}
			} else {
				value = fieldValue.Float()
			}
			if this.intParamAllowZero {
				if value != nil {
					strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
					this.args = append(this.args, value)
				}
			} else {
				if value != nil && !isEmpty(value) {
					strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
					this.args = append(this.args, value)
				}
			}

		case reflect.Struct:
			if objType.Field(i).Type.Name() == "Time" {
				value = fieldValue.Interface()
				if t, ok := value.(time.Time); ok {
					if t.IsZero() {
						value = nil
					} else {
						value = time.Unix(t.Unix(), 0).Format("2006-01-02 15:04:05")
					}
				}
				if value != nil && value != "" {

					strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
					this.args = append(this.args, value)
				}
			} else {
				log4go.Error("未处理的类型")
			}

		case reflect.Array:        //[]byte
			switch objType.Field(i).Type.Elem().Kind() {
			case reflect.Int8:
				if fieldValue.Kind() == reflect.Ptr {
					if fieldValue.IsNil() {
						value = nil
					} else {
						value = string(fieldValue.Elem().Bytes())
					}
				} else {
					value = string(fieldValue.Bytes())
				}
				if value != nil && value != "" {

					strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
					this.args = append(this.args, value)
				}
			default:
				log4go.Error("未处理的类型：%s", fieldType)

			}
		case reflect.Slice:        //[]byte
			switch objType.Field(i).Type.Elem().Kind() {
			case reflect.Uint8:
				if fieldValue.Kind() == reflect.Ptr {
					if fieldValue.IsNil() {
						value = nil
					} else {
						value = string(fieldValue.Elem().Bytes())
					}
				} else {
					value = fieldValue.Bytes()
				}
				if value != nil && value != "" {
					strsets = append(strsets, FieldFormat(fmt.Sprintf("%s", fieldName)))
					this.args = append(this.args, value)

				}
			default:
				log4go.Error("未处理的类型：%s", fieldType)

			}

		default:
			log4go.Error("类型转换错误：%s", fieldType)
		}
	}

	cols := strings.Join(strsets, ",")
	this.query = getInsertQuery(cols)
	this.isinsert = true
	return this
}

func (this*SQL)From(tableName string) *SQL {
	if len(this.from) == 0 {
		this.from = tableName
	} else {
		this.from = this.from + "," + tableName
	}
	return this
}

func (this*SQL)Where(paramWhere ... string) *SQL {
	this.where = where + query(paramWhere...)
	return this
}

func (this*SQL)ON(on ... string) *SQL {
	this.on = query(on...)
	return this
}

func (this*SQL)GroupBy(groupParam ... string) *SQL {
	this.group = group + strings.Join(groupParam, ",")
	return this
}

func (this*SQL)OrderBy(groupParam ... string) *SQL {
	this.order = order + strings.Join(groupParam, ",")
	return this
}

func (this*SQL)Args(args...interface{}) *SQL {
	this.args = args
	return this
}

//func In(params []interface{}) string {
//	strparams:=[]string{}
//	//for _,v:=range params{
//	//	strparams=append(strparams,fmt.Sprintf("%s",v))
//	//}
//	return " in"+LEFT+strings.Join(strparams,",")+RIGHT
//}
//
//func NotIn(params []interface{})  string{
//	strparams:=[]string{}
//	for _,v:=range params{
//		strparams=append(strparams,fmt.Sprintf("%s",v))
//	}
//	return " not in"+LEFT+strings.Join(strparams,",")+RIGHT
//}

func (this*SQL)GenerateCommand() string {
	if len(this.from) == 0 {
		log.Fatalf("GenerateCommand 没有目标的表名！ ")
	}
	var strSql string

	//参数分析
	newwhere := this.where
	limit := len(newwhere)
	limit = strings.LastIndex(newwhere[0:limit], PARAM)
	var pmax = len(this.args) - 1
	for {
		if limit < 0 {
			break
		} else {
			if rval := reflect.ValueOf(this.args[pmax]); rval.Kind() == reflect.Array || rval.Kind() == reflect.Slice {
				size := rval.Len()
				newwhere = newwhere[0:limit] + inparam[0:size * 2 - 1] + newwhere[limit + 1:]
			}
			pmax--
			limit--
			limit = strings.LastIndex(newwhere[0:limit], PARAM)
		}
	}

	if this.isselect {
		if len(this.query) == 0 {
			this.query = ALL
		}
		if (strings.Count(this.on, "") - 1) > 0 {
			strSql = selecT + this.query + from + this.from + on + this.on + newwhere + this.group + this.order
		} else {
			strSql = selecT + this.query + from + this.from + newwhere + this.group + this.order
		}

	} else if this.isupdate {
		if len(this.query) == 0 {
			log.Fatalf("update 没有操作！")
		}
		if len(newwhere) == 0 {
			log.Fatalf("update 没有条件")
		}
		strSql = update + this.from + set + this.query + newwhere
	} else if this.isdelete {
		if len(newwhere) == 0 {
			log.Fatalf("delete 没有条件")
		}
		strSql = delete + this.from + newwhere
	} else if this.isinsert {
		if len(this.query) == 0 {
			log.Fatalf("insert 没有操作！")
		}
		strSql = insert + this.from + this.query
	} else {
		log.Fatalf("GenerateCommand 没有选择操作类型")
	}

	//fmt.Println("GenerateCommand", strSql)
	return strSql
}

func (this*SQL) GetArgs() []interface{} {
	//处理in等数组问题
	//retargs:=make([]interface{},len(this.args))
	//for i,_:=range this.args{
	//	t:= reflect.ValueOf(this.args[i])
	//	if t.Kind() == reflect.Ptr{
	//		t=t.Elem()
	//	}
	//	if t.Kind()== reflect.Array || t.Kind()== reflect.Slice {
	//		arg:= this.args[i]
	//		vf:= reflect.ValueOf(arg)
	//		if vf.Len()>0{
	//			format:=""
	//			if vf.Index(0).Kind()== reflect.String{
	//				format=`"%s"`
	//			}else{
	//				format="%v"
	//			}
	//			vals:=make([]string,vf.Len())
	//			for l:=0;l<vf.Len();l++ {
	//				v:= fmt.Sprintf(format,vf.Index(l).Interface())
	//				vals[l]=v
	//			}
	//			retargs[i]=strings.Join(vals,",")
	//		}else{
	//			log4go.Error("参数值为空")
	//			retargs[i]=""
	//		}
	//	}else{
	//		retargs[i]=this.args[i]
	//	}
	//}
	//fmt.Println("GetArgs",retargs)
	//fmt.Println("GetArgs", this.args)
	return this.args
}

func SelectA(column string) string {
	if column == ALL || column == COUNT {
		return "a." + column
	} else {
		return "a." + FieldFormat(column)
	}

}

func SelectB(column string) string {
	if column == ALL || column == COUNT {
		return "b." + column
	} else {
		return "b." + FieldFormat(column)
	}
}

func SelectC(column string) string {
	if column == ALL || column == COUNT {
		return "c." + column
	} else {
		return "c." + FieldFormat(column)
	}
}

func SelectD(column string) string {
	if column == ALL || column == COUNT {
		return "d." + column
	} else {
		return "d." + FieldFormat(column)
	}
}

func Tables(table ...string) string {
	var tablename string = ""
	for i, t := range table {
		if len(tablename) > 0 {
			tablename += ","
		}
		tablename += t + " " + tablep[i:i+1]
	}
	return tablename
}

func CascadeTables(table ...string) string {
	var tablename string = ""
	for i, t := range table {
		if len(tablename) > 0 {
			tablename += LEFTJOIN
		}
		tablename += t + " " + tablep[i:i+1]
	}
	return tablename
}

func FieldFormat(field string) string {
	upfiled := strings.ToUpper(field);
	bup := []byte(upfiled)
	blower := []byte(strings.ToLower(field))
	bs := []byte(field)
	newfiled := []byte{}
	for i := 0; i < len(upfiled); i++ {
		if i == 0 {
			newfiled = append(newfiled, bs[i])
		} else if bup[i] == bs[i] && unicode.IsLetter(rune(bup[i])) {
			newfiled = append(newfiled, '_', bup[i])
		} else {
			newfiled = append(newfiled, blower[i])
		}
	}
	return string(newfiled)
}

