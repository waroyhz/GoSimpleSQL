package shell

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/astaxie/beego/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"bufio"
)

func command(strcmd ...string) string {
	//fmt.Println("exec command:",strcmd)
	cmd := exec.Command(strcmd[0], strcmd[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	//cmd.ExtraFiles = []*os.File{tf}

	if err := cmd.Run(); err != nil {
		log4go.Error(err, stderr.String())
	}

	result := stdout.String()
	return result

	//stdout, err := cmd.StdoutPipe()
	//defer stdout.Close()
	//if err != nil {
	//	fmt.Println("StdoutPipe: " + err.Error())
	//	return ""
	//}
	//
	//stderr, err := cmd.StderrPipe()
	//defer stderr.Close()
	//if err != nil {
	//	fmt.Println("StderrPipe: ", err.Error())
	//	return ""
	//}
	//
	//if err := cmd.Start(); err != nil {
	//	fmt.Println("Start: ", err.Error())
	//	return ""
	//}
	//
	//bytesErr, err := ioutil.ReadAll(stderr)
	//if err != nil {
	//	fmt.Println("ReadAll stderr: ", err.Error())
	//	return ""
	//}
	//
	//if len(bytesErr) != 0 {
	//	fmt.Printf("stderr is not nil: %s", bytesErr)
	//	return ""
	//}
	//
	//bytes, err := ioutil.ReadAll(stdout)
	//if err != nil {
	//	fmt.Println("ReadAll stdout: ", err.Error())
	//	return ""
	//}
	//
	//if err := cmd.Wait(); err != nil {
	//	fmt.Println("Wait: ", err.Error())
	//	return ""
	//}
	//
	////fmt.Println(bytes)
	////fmt.Printf("stdout: %s", bytes)
	//return string(bytes)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetCurrentPath() string {
	return strings.Replace(command("pwd"), "\n", "", -1)
}

// 拷贝文件夹
func CopyDir(src, dst string) bool {
	cmdrm := strings.Split(fmt.Sprintf("rm -rf %s", dst), " ")
	result := command(cmdrm...)
	cmdcp := strings.Split(fmt.Sprintf("cp -R %s %s", src, dst), " ")
	result += command(cmdcp...)
	return true
}

func MakeDir(path string) string  {
	cmdmk := strings.Split(fmt.Sprintf("mkdir -p %s", path), " ")
	result := command(cmdmk...)
	return result
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RunSMSClientSendSMS(phone, codestring, othermsg string) bool {
	file := fmt.Sprintf("%s%s", GetCurrentPath(), "/tools/SMSClient.jar")
	if utils.FileExists(file) {
		cmd := exec.Command("java", "-jar", file, "sendSMS", phone, codestring, othermsg)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		//cmd.ExtraFiles = []*os.File{tf}

		if err := cmd.Run(); err != nil {
			log4go.Error(err, stderr.String())
		}

		result := stdout.String()
		//fmt.Print("out:" ,result)
		lindex := strings.Index(result, "{")
		rindex := strings.Index(result, "}")
		if lindex > 0 && rindex > lindex {
			strjson := []byte(result)[lindex : rindex+1]
			fmt.Print("jsondata:", string(strjson))
			p := struct {
				Success bool `json:"success"`
			}{}
			json.Unmarshal(strjson, &p)

			return p.Success
		}
		log4go.Error(errors.New("发送短信失败！"), result)
	} else {
		log4go.Error(errors.New(file + " 文件不存在"))
		log.Println(errors.New(file + " 文件不存在"))
	}
	return false
}

//获取剩余短信容量
func RunSMSClientGetBalance() (float32, bool) {
	file := fmt.Sprintf("%s%s", GetCurrentPath(), "/tools/SMSClient.jar")
	if utils.FileExists(file) {
		cmd := exec.Command("java", "-jar", file, "getBalance")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		//cmd.ExtraFiles = []*os.File{tf}

		if err := cmd.Run(); err != nil {
			log4go.Error(err, stderr.String())
		}

		result := stdout.String()
		//fmt.Print("out:" ,result)
		lindex := strings.Index(result, "{")
		rindex := strings.Index(result, "}")
		if lindex > 0 && rindex > lindex {
			strjson := []byte(result)[lindex : rindex+1]
			fmt.Print("jsondata:", string(strjson))
			p := struct {
				Success bool    `json:"success"`
				Balance float32 `json:"balance"`
			}{}
			json.Unmarshal(strjson, &p)

			return p.Balance, p.Success
		}
		log4go.Error(errors.New("发送短信失败！"), result)
	} else {
		log4go.Error(errors.New(file + " 文件不存在"))
		log.Println(errors.New(file + " 文件不存在"))
	}
	return 0, false
}

func ReadFileString(filename string) (string,error)  {
	if f,e:= os.OpenFile(filename,os.O_RDONLY,0666);e==nil{
		defer f.Close()
		read:= bufio.NewReader(f)
		return read.ReadString(0)
	}else {
		return "",e
	}
}

func DeleteFile(filename string) error  {
	return os.Remove(filename)
}