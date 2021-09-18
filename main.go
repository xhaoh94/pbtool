package main

import (
	"flag"
	"fmt"
	"pbtool/conf"
	"pbtool/forms"
	"pbtool/parse"

	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
	_ "github.com/ying32/liblclbinres"
)

func main() {
	cfgId := *flag.Int("id", 0, "生产代码类型")
	flag.Parse()
	if conf.ReadCfg() {
		runApp(cfgId)
	}
}

func runApp(cfgId int) {

	if cfgId == 0 {
		vcl.Application.Initialize()
		vcl.Application.SetMainFormOnTaskBar(true)
		vcl.Application.CreateForm(&forms.MainForm)
		vcl.Application.CreateForm(&forms.CfgForm)
		vcl.Application.Run()
	} else {
		cfg, ok := conf.GetCfgs()[cfgId]
		if !ok {
			fmt.Printf("没有找到对应的配置ID:[%d]", cfgId)
			return
		}
		b := parse.Parse(cfg)
		if b {
			fmt.Println("导出成功")
		}
	}

	// vcl.RunApp(&forms.MainForm, &forms.CfgForm)
}
