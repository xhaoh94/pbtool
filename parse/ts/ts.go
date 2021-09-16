package ts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"pbtool/conf"
	"regexp"
	"strconv"
	"strings"

	"github.com/ying32/govcl/vcl"
)

var (
	messageReg      *regexp.Regexp
	messageTitleReg *regexp.Regexp
	enumReg         *regexp.Regexp
	enumTitleReg    *regexp.Regexp
	csReg           *regexp.Regexp
	scReg           *regexp.Regexp
	rpcReg          *regexp.Regexp
	contextReg      *regexp.Regexp

	fileName string
)

func init() {
	messageReg = regexp.MustCompile(`message ([^}]+)}`)
	messageTitleReg = regexp.MustCompile(`message ([^{]+)`)

	enumReg = regexp.MustCompile(`enum ([^}]+)}`)
	enumTitleReg = regexp.MustCompile(`enum ([^{]+)`)

	csReg = regexp.MustCompile(`cs=([\d]+)`)
	scReg = regexp.MustCompile(`sc=([\d]+)`)
	rpcReg = regexp.MustCompile(`rpc<([^>]+)>`)
	contextReg = regexp.MustCompile(`{([^}]+)}`)

}

func Parse(cfg *conf.OutCfg) {
	tsCfg := &conf.TsCfg{}
	err := json.Unmarshal([]byte(cfg.Context), tsCfg)
	if err != nil {
		vcl.ShowMessage("tsParse 解析Json报错")
		return
	}
	out := make([]string, 0)
	FilePathContent(cfg.InPath, &out)
	for _, v := range out {
		b := parseMessage(v)
		if !b {
			return
		}
		parseEnum(v)
		parseRPC(v)
	}
	write(cfg.OutPath, tsCfg.Ns, tsCfg.UseModule, tsCfg.CreateJson, tsCfg.FileName)
	if tsCfg.CreateJson {
		writeJSON(tsCfg.OutJsonPath + tsCfg.JsonName)
	}
}

func parseRPC(str string) {
	strs := rpcReg.FindAllString(str, -1)
	for _, context := range strs {
		if strings.Index(context, ":") < 0 {
			continue
		}
		s := &RpcStruct{}
		rpcMatched := rpcReg.FindStringSubmatch(context)
		if len(rpcMatched) == 2 {
			str := strings.Replace(rpcMatched[1], " ", "", -1)
			str = strings.Replace(str, "<", "", -1)
			str = strings.Replace(str, ">", "", -1)
			ss := strings.Split(str, ":")
			s.Req = ss[0]
			s.Rsp = ss[1]
		}
		Rpcs = append(Rpcs, s)
	}
}
func parseEnum(str string) {
	strs := enumReg.FindAllString(str, -1)
	for _, context := range strs {
		s := &EnumStruct{}
		s.Datas = make([][]string, 0)
		titleMatched := enumTitleReg.FindStringSubmatch(context)
		if len(titleMatched) == 2 {
			s.Title = titleMatched[1]
		}

		contMatched := contextReg.FindStringSubmatch(context)
		if len(contMatched) == 2 {
			lines := strings.Split(contMatched[1], "\n")
			for i := 0; i < len(lines); i++ {
				line := lines[i]
				if strings.Index(line, "=") < 0 {
					continue
				}
				c := strings.Split(line, "=")
				if strings.Index(c[0], "//") >= 0 {
					continue
				}
				datas := make([]string, 0)

				c[1] = strings.Replace(c[1], " ", "", -1)
				endIndex := strings.Index(c[1], ";")
				tag := c[1][0:endIndex]
				datas = append(datas, tag)
				ts := strings.Split(c[0], " ")
				for _, v := range ts {
					if v != "" {
						datas = append(datas, v)
					}
				}
				s.Datas = append(s.Datas, datas)
			}
		}
		IsEnum[s.Title] = true
		Enums = append(Enums, s)
	}
}

func parseMessage(str string) bool {
	strs := messageReg.FindAllString(str, -1)
	for _, context := range strs {
		s := &MessageStruct{}
		s.Datas = make([][]string, 0)
		csMatched := csReg.FindStringSubmatch(context)
		if len(csMatched) == 2 {
			num, err := strconv.Atoi(csMatched[1])
			if err != nil {
				fmt.Printf("cmd非int类型%s", csMatched[1])
				return false
			}
			s.Cs = uint32(num)
		} else {
			scMatched := scReg.FindStringSubmatch(context)
			if len(scMatched) == 2 {
				num, err := strconv.Atoi(scMatched[1])
				if err != nil {
					fmt.Printf("cmd非int类型%s", scMatched[1])
					return false
				}
				s.Sc = uint32(num)
			}
		}

		titleMatched := messageTitleReg.FindStringSubmatch(context)
		if len(titleMatched) == 2 {
			s.Title = titleMatched[1]
		}
		contMatched := contextReg.FindStringSubmatch(context)
		if len(contMatched) == 2 {
			var startindex int
			if s.Cs > 0 || s.Sc > 0 {
				startindex = 1
			}
			lines := strings.Split(contMatched[1], "\n")
			for i := startindex; i < len(lines); i++ {
				line := lines[i]
				if strings.Index(line, "=") < 0 {
					continue
				}
				c := strings.Split(line, "=")
				if strings.Index(c[0], "//") >= 0 {
					continue
				}

				datas := make([]string, 0)
				c[1] = strings.Replace(c[1], " ", "", -1)
				endIndex := strings.Index(c[1], ";")
				tag := c[1][0:endIndex]
				datas = append(datas, tag)

				ts := strings.Split(c[0], " ")
				for _, v := range ts {
					if v != "" && v != "repeated" {
						datas = append(datas, v)
					}
				}
				if strings.Index(c[0], "repeated") >= 0 {
					datas = append(datas, "1")
				} else {
					datas = append(datas, "0")
				}
				s.Datas = append(s.Datas, datas)
			}
		}
		Messages = append(Messages, s)
	}
	return true
}

//write 写入
func write(OutPath string, NameSpace string, UseModule bool, CreateJson bool, FileName string) {
	_, err := os.Stat(OutPath)
	if err != nil {
		os.Mkdir(OutPath, 0777)
	}

	str := "namespace " + NameSpace + "{\n"
	if UseModule {
		str = Title + "export " + str
	} else {
		str = Title + str
	}
	str += writeCmd() + "\n"
	if !CreateJson {
		str += writeConf() + "\n"
	}
	for _, v := range Enums {
		str += "\texport const enum " + v.Title + " {\n"
		f := true
		for _, c := range v.Datas {
			if f {
				f = false
				str += "\t\t" + c[1] + "=" + c[0]
			} else {
				str += ",\n\t\t" + c[1] + "=" + c[0]
			}
		}
		str += "\n\t}\n"
	}
	for _, v := range Messages {
		str += "\texport interface " + v.Title + " {\n"
		for _, c := range v.Datas { //tag type name isArray
			isArray := c[len(c)-1] == "1"
			if isArray {
				str += "\t\t" + c[2] + "?:" + GetType(c[1]) + "[];\n"
			} else {
				str += "\t\t" + c[2] + "?:" + GetType(c[1]) + ";\n"
			}
		}
		str += "\t}\n"
	}
	str += "}"

	var d = []byte(str)
	err = ioutil.WriteFile(OutPath+FileName, d, 0666)
	if err != nil {
		fmt.Println("write ts fail")
	} else {
		fmt.Println("write ts success")
	}
}

func writeCmd() string {

	strCs := "\texport const enum CS" + " {\n"
	strSc := "\texport const enum SC" + " {\n"
	for _, v := range Messages {
		if v.Cs > 0 {
			strCs += "\t\t" + v.Title + " = " + strconv.Itoa(int(v.Cs)) + ",\n"
		} else if v.Sc > 0 {
			strSc += "\t\t" + v.Title + " = " + strconv.Itoa(int(v.Sc)) + ",\n"
		}
	}
	strCs += "\t}\n"
	strSc += "\t}"

	return strCs + strSc
}

func writeConf() string {

	rpc := "\texport const rpcMap:{ [key: string]: string }={\n"
	f := true
	for k := 0; k < len(Rpcs); k++ {
		v := Rpcs[k]
		if f {
			f = false
			rpc += "\t\t" + GetString(v.Req) + ":" + GetString(v.Rsp)
		} else {
			rpc += ",\n\t\t" + GetString(v.Req) + ":" + GetString(v.Rsp)
		}
	}
	rpc += "\n\t}\n"
	cs := "\texport const csMap:{ [key: number]: string }={\n"
	sc := "\texport const csMap:{ [key: number]: string }={\n"
	message := "\texport const messageMap:{ [key: string]: string[][] }={\n"
	fcs := true
	fsc := true
	for j := 0; j < len(Messages); j++ {
		v := Messages[j]
		if v.Cs > 0 {
			if fcs {
				fcs = false
				cs += "\t\t" + strconv.Itoa(int(v.Cs)) + ":" + GetString(v.Title)
			} else {
				cs += ",\n\t\t" + strconv.Itoa(int(v.Cs)) + ":" + GetString(v.Title)
			}
		} else {
			if v.Sc > 0 {
				if fsc {
					fsc = false
					sc += "\t\t" + strconv.Itoa(int(v.Sc)) + ":" + GetString(v.Title)
				} else {
					sc += ",\n\t\t" + strconv.Itoa(int(v.Sc)) + ":" + GetString(v.Title)
				}
			}
		}

		message += "\t\t" + GetString(v.Title) + ":["
		for i := 0; i < len(v.Datas); i++ {
			c := v.Datas[i]
			message += "[" + GetString(c[0]) + "," + GetString(c[2]) + "," + GetId(c[1])
			isArray := c[len(c)-1] == "1"
			if isArray {
				message += "," + GetString("1")
			}
			if i == len(v.Datas)-1 {
				message += "]"
			} else {
				message += "],"
			}
		}
		if j == len(Messages)-1 {
			message += "]\n"
		} else {
			message += "],\n"
		}
	}
	message += "\t}\n"
	sc += "\n\t}\n"
	cs += "\n\t}\n"

	r := cs + sc
	if len(Rpcs) > 0 {
		r += rpc
	}
	r += message
	return r

}

func writeJSON(fileName string) {
	str := "{\n"

	rpcMap := "\t" + GetString("rpcMap") + ":{\n"
	f := true
	for k := 0; k < len(Rpcs); k++ {
		v := Rpcs[k]
		if f {
			f = false
			rpcMap += "\t\t" + GetString(v.Req) + ":" + GetString(v.Rsp)
		} else {
			rpcMap += ",\n\t\t" + GetString(v.Req) + ":" + GetString(v.Rsp)
		}
	}
	rpcMap += "\n\t},\n"

	csMap := "\t" + GetString("csMap") + ":{\n"
	scMap := "\t" + GetString("scMap") + ":{\n"
	messageMap := "\t" + GetString("messageMap") + ":{\n"
	fcs := true
	fsc := true
	for j := 0; j < len(Messages); j++ {
		v := Messages[j]
		if v.Cs > 0 {
			if fcs {
				fcs = false
				csMap += "\t\t" + GetString(strconv.Itoa(int(v.Cs))) + ":" + GetString(v.Title)
			} else {
				csMap += ",\n\t\t" + GetString(strconv.Itoa(int(v.Cs))) + ":" + GetString(v.Title)
			}
		} else if v.Sc > 0 {
			if fsc {
				fsc = false
				scMap += "\t\t" + GetString(strconv.Itoa(int(v.Sc))) + ":" + GetString(v.Title)
			} else {
				scMap += ",\n\t\t" + GetString(strconv.Itoa(int(v.Sc))) + ":" + GetString(v.Title)
			}
		}
		messageMap += "\t\t" + GetString(v.Title) + ":["
		for i := 0; i < len(v.Datas); i++ {
			c := v.Datas[i]
			messageMap += "[" + GetString(c[0]) + "," + GetString(c[2]) + "," + GetId(c[1])
			isArray := c[len(c)-1] == "1"
			if isArray {
				messageMap += "," + GetString("1")
			}
			if i == len(v.Datas)-1 {
				messageMap += "]"
			} else {
				messageMap += "],"
			}
		}
		if j == len(Messages)-1 {
			messageMap += "]\n"
		} else {
			messageMap += "],\n"
		}
	}
	messageMap += "\t}\n"
	csMap += "\n\t},\n"
	scMap += "\n\t},\n"
	str += csMap
	str += scMap
	str += rpcMap
	str += messageMap
	str += "}"

	var d = []byte(str)
	err := ioutil.WriteFile(fileName, d, 0666)
	if err != nil {
		fmt.Println("write json fail")
	}
	fmt.Println("write json success")
}
