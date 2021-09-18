package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	conf                *ConfStruct
	TypeComboxs         []string = []string{TC_TypeScript, TC_Golang, TC_CSharp, TC_Protobufjs}
	PBJS_Wrap_Comboxs   []string = []string{PBJS_Commonjs, PBJS_Deafult, PBJS_Amd, PBJS_Es6, PBJS_Closure}
	PBJS_Target_Comboxs []string = []string{PBJS_Staitc_Module, PBJS_Json, PBJS_Json_Module, PBJS_Proto2, PBJS_Proto3, PBJS_Static}
)

const (
	TC_TypeScript string = "TypeScript"
	TC_Golang     string = "Golang"
	TC_CSharp     string = "CSharp"
	TC_Protobufjs string = "Protobufjs"
)
const (
	PBJS_Deafult  string = "default"
	PBJS_Commonjs string = "commonjs"
	PBJS_Amd      string = "amd"
	PBJS_Es6      string = "es6"
	PBJS_Closure  string = "closure"

	PBJS_Json          string = "json"
	PBJS_Json_Module   string = "json-module"
	PBJS_Proto2        string = "proto2"
	PBJS_Proto3        string = "proto3"
	PBJS_Static        string = "static"
	PBJS_Staitc_Module string = "static-module"
)

type (
	ConfStruct struct {
		ProtocPath string
		Cfgs       map[int]*OutCfg
	}

	OutCfg struct {
		ID      int
		Name    string
		InPath  string
		OutPath string
		TagType string
		Context string
	}

	TsCfg struct {
		Ns          string
		CreateJson  bool
		UseModule   bool
		OutJsonPath string
		FileName    string
		JsonName    string
	}
	GoCfg struct {
		CreateCmd      bool
		Ns             string
		FileName       string
		ProtoGenGoPath string
	}
	CSharpCfg struct {
		CreateCmd bool
		Ns        string
		FileName  string
	}
	PbJsCfg struct {
		FileName  string
		CreateDts bool
		Target    string
		Wrap      string
		UseEs6    bool
	}
)

func initCfg() {
	conf = &ConfStruct{
		ProtocPath: "",
		Cfgs:       make(map[int]*OutCfg),
	}
	WriteCfg()
}
func SetProtoPath(p string) {
	if conf.ProtocPath != p {
		conf.ProtocPath = p
		WriteCfg()
	}
}
func GetProtoPath() string {
	return conf.ProtocPath
}
func AddCfg(cfg *OutCfg) {
	conf.Cfgs[cfg.ID] = cfg
	WriteCfg()
}
func DelCfg(id int) {
	delete(conf.Cfgs, id)
	WriteCfg()
}
func GetCfgs() map[int]*OutCfg {
	return conf.Cfgs
}
func ReadCfg() bool {
	fileByte, fileErr := ioutil.ReadFile("pbtool.conf")
	if fileErr != nil {
		initCfg()
		return true
	}
	fileErr = json.Unmarshal(fileByte, &conf)
	if fileErr != nil {
		fmt.Println("读取config.json有误")
		return false
	}
	return true
}

func WriteCfg() {
	bytes, _ := json.Marshal(&conf)
	ioutil.WriteFile("pbtool.conf", bytes, 0777)
}
