package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	conf        *ConfStruct
	TypeComboxs []string = []string{TC_TypeScript}
)

const (
	TC_TypeScript string = "TypeScript"
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
