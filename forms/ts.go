package forms

import (
	"encoding/json"
	"path/filepath"
	"pbtool/conf"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func (f *TMainForm) tsPanel(tab *vcl.TTabSheet) {
	cfgs := conf.GetCfgs()
	cfg := cfgs[tab.Tag()]
	// cfg.Context
	tsCfg := &conf.TsCfg{}
	err := json.Unmarshal([]byte(cfg.Context), tsCfg)
	if err != nil {
		vcl.ShowMessage("tsPanel 解析Json报错")
		return
	}

	p := vcl.NewFrame(tab)
	p.SetParent(tab)
	p.SetAlign(types.AlClient)

	lb := vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("配置类型：")
	lb.SetTop(15)
	lb.SetLeft(10)

	typeEdit := vcl.NewEdit(p)
	typeEdit.SetParent(p)
	typeEdit.SetEnabled(false)
	typeEdit.SetBounds(100, 10, 200, 50)
	typeEdit.SetText(cfg.TagType)

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("proto文件目录：")
	lb.SetTop(45)
	lb.SetLeft(10)

	dirEdit := vcl.NewEdit(p)
	dirEdit.SetEnabled(false)
	dirEdit.SetParent(p)
	dirEdit.SetBounds(100, 40, 500, 50)
	dirEdit.SetText(cfg.InPath)

	btnDir := vcl.NewButton(p)
	btnDir.SetParent(p)
	btnDir.SetCaption("···")
	btnDir.SetBounds(620, 40, 80, 25)
	btnDir.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewSelectDirectoryDialog(nil)
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			dirEdit.SetText(path)
			if cfg.InPath != path {
				cfg.InPath = path
				updCfg(cfg, tsCfg)
			}
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出文件目录：")
	lb.SetTop(75)
	lb.SetLeft(10)

	outEdit := vcl.NewEdit(p)
	outEdit.SetEnabled(false)
	outEdit.SetParent(p)
	outEdit.SetBounds(100, 70, 500, 50)
	outEdit.SetText(cfg.OutPath)

	btnOut := vcl.NewButton(p)
	btnOut.SetParent(p)
	btnOut.SetCaption("···")
	btnOut.SetBounds(620, 70, 80, 25)
	btnOut.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewSelectDirectoryDialog(nil)
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			outEdit.SetText(path)
			if cfg.OutPath != path {
				cfg.OutPath = path
				updCfg(cfg, tsCfg)
			}
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出Json目录：")
	lb.SetTop(105)
	lb.SetLeft(10)

	jsonEdit := vcl.NewEdit(p)
	jsonEdit.SetEnabled(false)
	jsonEdit.SetParent(p)
	jsonEdit.SetBounds(100, 100, 500, 50)
	jsonEdit.SetText(tsCfg.OutJsonPath)

	btnJson := vcl.NewButton(p)
	btnJson.SetParent(p)
	btnJson.SetCaption("···")
	btnJson.SetBounds(620, 100, 80, 25)
	btnJson.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewSelectDirectoryDialog(nil)
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			if tsCfg.OutJsonPath != path {
				jsonEdit.SetText(path)
				tsCfg.OutJsonPath = path
				updCfg(cfg, tsCfg)
			}
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("命名空间:")
	lb.SetTop(135)
	lb.SetLeft(10)

	edit := vcl.NewEdit(p)
	edit.SetParent(p)
	edit.SetBounds(100, 130, 200, 50)
	edit.SetOnExit(func(sender vcl.IObject) {
		if tsCfg.Ns != edit.Text() {
			tsCfg.Ns = edit.Text()
			updCfg(cfg, tsCfg)
		}
	})
	edit.SetText(tsCfg.Ns)

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出文件名:")
	lb.SetTop(165)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(p)
	fileEdit.SetParent(p)
	fileEdit.SetBounds(100, 160, 200, 50)
	fileEdit.SetText(tsCfg.FileName)
	fileEdit.SetOnExit(func(sender vcl.IObject) {
		if tsCfg.FileName != fileEdit.Text() {
			tsCfg.FileName = fileEdit.Text()
			updCfg(cfg, tsCfg)
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出Json名:")
	lb.SetTop(195)
	lb.SetLeft(10)

	jsonNameEdit := vcl.NewEdit(p)
	jsonNameEdit.SetParent(p)
	jsonNameEdit.SetBounds(100, 190, 200, 50)
	jsonNameEdit.SetText(tsCfg.JsonName)
	jsonNameEdit.SetOnExit(func(sender vcl.IObject) {
		if tsCfg.JsonName != jsonNameEdit.Text() {
			tsCfg.JsonName = jsonNameEdit.Text()
			updCfg(cfg, tsCfg)
		}
	})

	cb1 := vcl.NewCheckBox(p)
	cb1.SetParent(p)
	cb1.SetCaption("是否生成Json配置")
	cb1.SetTop(220)
	cb1.SetLeft(10)
	cb1.SetChecked(tsCfg.CreateJson)
	cb1.SetOnChange(func(sender vcl.IObject) {
		if tsCfg.CreateJson != cb1.Checked() {
			tsCfg.CreateJson = cb1.Checked()
			updCfg(cfg, tsCfg)
		}
	})

	cb2 := vcl.NewCheckBox(p)
	cb2.SetParent(p)
	cb2.SetCaption("是否增加export关键字")
	cb2.SetTop(240)
	cb2.SetLeft(10)
	cb2.SetChecked(tsCfg.UseModule)
	cb2.SetOnChange(func(sender vcl.IObject) {
		if tsCfg.UseModule != cb2.Checked() {
			tsCfg.UseModule = cb2.Checked()
			updCfg(cfg, tsCfg)
		}
	})

}

func (f *TCfgForm) tsPanel() *PnlCfg {
	pnlCfg := &PnlCfg{
		Cfg: &conf.TsCfg{
			Ns:          "proto",
			CreateJson:  true,
			UseModule:   true,
			OutJsonPath: "",
			FileName:    "DymPbCode.ts",
			JsonName:    "DymPbCfg.json",
		},
		Pnl: vcl.NewFrame(f.pContext),
	}

	// pnlCfg.Pnl.SetParent(f.pContext)
	pnlCfg.Pnl.SetAlign(types.AlClient)

	lb := vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("命名空间:")
	lb.SetTop(15)
	lb.SetLeft(10)

	edit := vcl.NewEdit(pnlCfg.Pnl)
	edit.SetParent(pnlCfg.Pnl)
	edit.SetBounds(100, 10, 200, 50)
	edit.SetText(pnlCfg.Cfg.(*conf.TsCfg).Ns)
	edit.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.TsCfg).Ns = edit.Text()
	})

	lb = vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("导出文件名:")
	lb.SetTop(45)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(pnlCfg.Pnl)
	fileEdit.SetParent(pnlCfg.Pnl)
	fileEdit.SetBounds(100, 40, 200, 50)
	fileEdit.SetText(pnlCfg.Cfg.(*conf.TsCfg).FileName)
	fileEdit.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.TsCfg).FileName = fileEdit.Text()
	})

	lb = vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("导出Json名:")
	lb.SetTop(75)
	lb.SetLeft(10)

	jsonNameEdit := vcl.NewEdit(pnlCfg.Pnl)
	jsonNameEdit.SetParent(pnlCfg.Pnl)
	jsonNameEdit.SetBounds(100, 70, 200, 50)
	jsonNameEdit.SetText(pnlCfg.Cfg.(*conf.TsCfg).JsonName)
	jsonNameEdit.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.TsCfg).JsonName = jsonNameEdit.Text()
	})

	lb = vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("导出Json目录：")
	lb.SetTop(105)
	lb.SetLeft(10)

	jsonEdit := vcl.NewEdit(pnlCfg.Pnl)
	jsonEdit.SetEnabled(false)
	jsonEdit.SetParent(pnlCfg.Pnl)
	jsonEdit.SetBounds(100, 100, 400, 50)

	btnJson := vcl.NewButton(pnlCfg.Pnl)
	btnJson.SetParent(pnlCfg.Pnl)
	btnJson.SetCaption("···")
	btnJson.SetBounds(510, 100, 80, 25)
	btnJson.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewSelectDirectoryDialog(nil)
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			jsonEdit.SetText(path)
			pnlCfg.Cfg.(*conf.TsCfg).OutJsonPath = path
		}
	})

	cb1 := vcl.NewCheckBox(pnlCfg.Pnl)
	cb1.SetParent(pnlCfg.Pnl)
	cb1.SetCaption("是否生成Json配置")
	cb1.SetTop(135)
	cb1.SetLeft(10)
	cb1.SetChecked(pnlCfg.Cfg.(*conf.TsCfg).CreateJson)
	cb1.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.TsCfg).CreateJson = cb1.Checked()
	})

	cb2 := vcl.NewCheckBox(pnlCfg.Pnl)
	cb2.SetParent(pnlCfg.Pnl)
	cb2.SetCaption("是否增加export关键字")
	cb2.SetTop(155)
	cb2.SetLeft(10)
	cb2.SetChecked(pnlCfg.Cfg.(*conf.TsCfg).UseModule)
	cb2.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.TsCfg).UseModule = cb2.Checked()
	})

	return pnlCfg
}
