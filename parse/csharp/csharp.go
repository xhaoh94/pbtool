package csharp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"pbtool/common"
	"pbtool/conf"
	"regexp"
	"strconv"
	"strings"

	"github.com/ying32/govcl/vcl"
)

var (
	messageReg      *regexp.Regexp
	messageTitleReg *regexp.Regexp
	csReg           *regexp.Regexp
	scReg           *regexp.Regexp
	contextReg      *regexp.Regexp

	fileName string
)

func init() {
	messageReg = regexp.MustCompile(`message ([^}]+)}`)
	messageTitleReg = regexp.MustCompile(`message ([^{]+)`)

	csReg = regexp.MustCompile(`cs=([\d]+)`)
	scReg = regexp.MustCompile(`sc=([\d]+)`)
	contextReg = regexp.MustCompile(`{([^}]+)}`)

}

func Parse(cfg *conf.OutCfg) bool {
	protocPath := conf.GetProtoPath()
	if protocPath == "" {
		vcl.ShowMessage("protoc路径为空")
		return false
	}
	goCfg := &conf.GoCfg{}
	err := json.Unmarshal([]byte(cfg.Context), goCfg)
	if err != nil {
		vcl.ShowMessage("csharpParse 解析Json报错")
		return false
	}

	files := make([]string, 0)
	common.Files(cfg.InPath, &files)

	str := protocPath + " -I=" + cfg.InPath + " "
	str += "--csharp_out=" + cfg.OutPath + " "
	for _, v := range files {
		exe := str + v
		b := common.RunExe(exe, "start")
		if !b {
			log.Println("csharpParse 解析协议出错：" + str)
		}
	}

	if goCfg.CreateCmd {
		out := make([]string, 0)
		common.FilePathContent(cfg.InPath, &out)
		for _, v := range out {
			b := parseMessage(v)
			if !b {
				return false
			}
		}
		write(cfg.OutPath, goCfg.Ns, goCfg.FileName)
	}
	return true
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

func write(OutPath string, NameSpace string, FileName string) {
	_, err := os.Stat(OutPath)
	if err != nil {
		os.Mkdir(OutPath, 0777)
	}

	strCs := "public enum CS {\n"
	strSc := "public enum SC {\n"
	for _, v := range Messages {
		if v.Cs > 0 {
			if NameSpace != "" {
				strCs += "\t"
			}
			strCs += "\t" + v.Title + " = " + strconv.Itoa(int(v.Cs)) + "\n"
		}
		if v.Sc > 0 {
			if NameSpace != "" {
				strSc += "\t"
			}
			strSc += "\t" + v.Title + " = " + strconv.Itoa(int(v.Sc)) + "\n"
		}
	}
	if NameSpace != "" {
		strCs += "\t"
	}
	strCs += "}"
	if NameSpace != "" {
		strSc += "\t"
	}
	strSc += "}"
	str := common.Title
	if NameSpace != "" {
		str += "namespace " + NameSpace + "{%s\n}"
		strCs = "\n\t" + strCs
		strSc = "\n\t" + strSc
	} else {
		str += "%s"
		strCs = "\n" + strCs
		strSc = "\n" + strSc
	}
	context := fmt.Sprintf(str, strCs+strSc)
	var d = []byte(context)
	err = ioutil.WriteFile(OutPath+"/"+FileName, d, 0666)
	if err != nil {
		fmt.Println("write golang fail")
	} else {
		fmt.Println("write golang success")
	}
}
