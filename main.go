package main

import (
	"pbtool/conf"
	"pbtool/forms"

	_ "github.com/ying32/govcl/pkgs/winappres"
	"github.com/ying32/govcl/vcl"
)

func main() {
	if conf.ReadCfg() {
		runApp()
	}
}

func runApp() {

	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&forms.MainForm)
	vcl.Application.CreateForm(&forms.CfgForm)
	vcl.Application.Run()
	// vcl.RunApp(&forms.MainForm, &forms.CfgForm)
}
