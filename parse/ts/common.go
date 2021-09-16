package ts

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type (
	MessageStruct struct {
		Cs    uint32
		Sc    uint32
		Title string
		Datas [][]string
	}
	EnumStruct struct {
		Title string
		Datas [][]string
	}
	RpcStruct struct {
		Req string
		Rsp string
	}
)

var (
	Title    string = "// Generated by https://github.com/xhaoh94/protogen\n// DO NOT EDIT!\n"
	Messages []*MessageStruct
	Enums    []*EnumStruct
	Rpcs     []*RpcStruct
	IsEnum   map[string]bool = make(map[string]bool)
)

func GetString(str string) string {
	return "\"" + str + "\""
}
func GetType(str string) string {
	switch str {
	case "bool":
		return "boolean"
	case "string":
		return "string"
	case "bytes":
		return "Uint8Array"
	case "float":
		return "number"
	case "double":
		return "number"
	case "enum":
		return "enum"
	case "int32":
		return "number"
	case "int64":
		return "number|Long"
	case "uint32":
		return "number"
	case "uint64":
		return "number|Long"
	case "sint32":
		return "number"
	case "sint64":
		return "number|Long"
	case "fixed32":
		return "number"
	case "fixed64":
		return "number|Long"
	case "sfixed32":
		return "number"
	case "sfixed64":
		return "number|Long"
	default:
		return str
	}
}
func GetId(str string) string {
	return GetString(cov(str))
}
func cov(str string) string {
	if IsEnum[str] {
		return "5"
	}
	switch str {
	case "bool":
		return "0"
	case "string":
		return "1"
	case "bytes":
		return "2"
	case "float":
		return "3"
	case "double":
		return "4"
	case "int32":
		return "6"
	case "int64":
		return "7"
	case "uint32":
		return "8"
	case "uint64":
		return "9"
	case "sint32":
		return "10"
	case "sint64":
		return "11"
	case "fixed32":
		return "12"
	case "fixed64":
		return "13"
	case "sfixed32":
		return "14"
	case "sfixed64":
		return "15"
	default:
		return str
	}
}
func FilePathContent(path string, out *[]string) {
	if string(path[len(path)-1]) != "\\" {
		path += "\\"
	}
	fs, _ := ioutil.ReadDir(path)
	for _, file := range fs {
		if file.IsDir() {
			FilePathContent(path+file.Name()+"\\", out)
		} else {
			if strings.Index(file.Name(), ".") < 0 {
				FilePathContent(path+file.Name()+"\\", out)
				continue
			}
			if strings.Index(file.Name(), "proto") < 0 {

				continue
			}
			fmt.Println(file.Name())
			f, err := ioutil.ReadFile(path + file.Name())
			if err != nil {
				fmt.Println("read fail", err)
			}
			*out = append(*out, string(f))
		}
	}
}