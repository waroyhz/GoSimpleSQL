package sql

import (
	"testing"
)

func Test_NewCommand(t *testing.T) {
	sql:= NewCommand("test")
	if sql.from!="test" { //try a unit test on function
		t.Error(t.Name(),sql.from) // 如果不是如预期的那么就报错
	}
}

func Test_Select(t *testing.T) {
	sql:= NewCommand("test").Select("a1")
	if sql.query!="a1" { //try a unit test on function
		t.Error(t.Name(),sql.query) // 如果不是如预期的那么就报错
	}
}

func Test_Delete(t *testing.T) {
	sql:= NewCommand("test").Select("a1")
	sql= sql.Delete(WhereEQ{"a2","a3"})
	if !sql.isdelete{
		t.Error(t.Name(),sql.isdelete)
	}
	if sql.where!=" where a2=a3" {
		t.Error(t.Name(),sql.where)
	}
}

func Test_Deletes(t *testing.T) {
	sql:= NewCommand("test").Select("a1")
	sql= sql.Deletes()
	if !sql.isdelete{
		t.Error(t.Name(),sql.isdelete)
	}

}

func Test_Update(t *testing.T) {
	sql:= NewCommand("test").Select("a1")
	sql= sql.Update(Set{"a1","a2"})
	if !sql.isupdate{
		t.Error(t.Name(),sql.isdelete)
	}
	if sql.query!="a1=a2" {
		t.Error(t.Name(),sql.query)
	}
}

func Test_IntParamAllowZero(t*testing.T){
	sql:= NewCommand("test").Select("a1")
	sql= sql.IntParamAllowZero()
	if !sql.intParamAllowZero{
		t.Error(t.Name(),sql.intParamAllowZero)
	}
}

func Test_Insert(t*testing.T){
	var entity = struct {
		Str string
		Int int
		Float float32
		Bool bool
	}{"str",1,1.1,true}

	sql:= NewCommand("test")
	sql= sql.Insert(&entity)
	if sql.query!=" (Str,Int,Float,Bool)  values(?,?,?,?) "{
		t.Error(t.Name(),sql.query)
	}
	f3:=float32(sql.args[2].(float64))
	if !(sql.args[0].(string)=="str" && sql.args[1].(int64)==1 && f3== 1.1 && sql.args[3].(bool)==true){
		t.Error(t.Name(),sql.args[0],sql.args[1],f3,sql.args[3])
	}
}


func TestSQL_From(t*testing.T){
	sql:= NewCommand("test").From("a1")
	if sql.from!="test,a1"{
		t.Error(t.Name(),sql.from)
	}
}

func TestSQL_Where(t*testing.T){
	sql:= NewCommand("test").Where("a1")
	if sql.where!=" where a1"{
		t.Error(t.Name(),sql.where)
	}
}

func TestSQL_ON(t*testing.T){
	sql:= NewCommand("test").ON("a1","=","a2")
	if sql.on!="a1=a2"{
		t.Error(t.Name(),sql.on)
	}
}

func TestSQL_GroupBy(t*testing.T){
	sql:= NewCommand("test").GroupBy("a1","a2")
	if sql.group!=" group by a1,a2"{
		t.Error(t.Name(),sql.group)
	}
}

func TestSQL_OrderBy(t*testing.T){
	sql:= NewCommand("test").OrderBy("a1","a2")
	if sql.order!=" order by a1,a2"{
		t.Error(t.Name(),sql.order)
	}
}


func TestSQL_Args(t*testing.T){
	sql:= NewCommand("test").Args("a1","a2")
	if !(sql.args[0]=="a1" && sql.args[1]=="a2"){
		t.Error(t.Name(),sql.args[0],sql.args[1])
	}
}

func TestSQL_GetArgs(t*testing.T){
	sql:= NewCommand("test").Args("a1","a2")
	if !(sql.args[0]=="a1" && sql.args[1]=="a2"){
		t.Error(t.Name(),sql.args[0],sql.args[1])
	}
}

func Test_SelectA(t*testing.T){
	if !(SelectA("a1") =="a.a1"){
		t.Error(t.Name(),SelectA("a1"))
	}
}

func Test_SelectB(t*testing.T){
	if !(SelectB("a1") =="b.a1"){
		t.Error(t.Name(),SelectB("a1"))
	}
}

func Test_SelectC(t*testing.T){
	if !(SelectC("a1") =="c.a1"){
		t.Error(t.Name(),SelectC("a1"))
	}
}

func Test_SelectD(t*testing.T){
	if !(SelectD("a1") =="d.a1"){
		t.Error(t.Name(),SelectD("a1"))
	}
}

func Test_Tables(t*testing.T){
	if !(Tables("a1","a2") =="a1 a,a2 b"){
		t.Error(t.Name(),Tables("a1","a2"))
	}
}

func Test_CascadeTables(t*testing.T){
	if !(CascadeTables("a1","a2") =="a1 a left join a2 b"){
		t.Error(t.Name(),CascadeTables("a1","a2"))
	}
}

//func Test_FieldFormat(t*testing.T){
//	if !(FieldFormat("aA1") =="a_A1"){
//		t.Error(t.Name(),FieldFormat("aA1"))
//	}
//	if !(FieldFormat("a1") =="a1"){
//		t.Error(t.Name(),FieldFormat("a1"))
//	}
//}