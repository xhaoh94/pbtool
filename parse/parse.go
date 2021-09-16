package parse

import (
	"pbtool/conf"
	"pbtool/parse/ts"
)

func Parse(cfg *conf.OutCfg) {
	switch cfg.TagType {
	case conf.TC_TypeScript:
		ts.Parse(cfg)
		break
	}
}
