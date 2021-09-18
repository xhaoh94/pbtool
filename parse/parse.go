package parse

import (
	"pbtool/conf"
	"pbtool/parse/csharp"
	"pbtool/parse/golang"
	"pbtool/parse/pbjs"
	"pbtool/parse/ts"
)

func Parse(cfg *conf.OutCfg) bool {
	var b bool
	switch cfg.TagType {
	case conf.TC_TypeScript:
		b = ts.Parse(cfg)
		break
	case conf.TC_Golang:
		b = golang.Parse(cfg)
		break
	case conf.TC_CSharp:
		b = csharp.Parse(cfg)
		break
	case conf.TC_Protobufjs:
		b = pbjs.Parse(cfg)
		break
	}
	return b

}
