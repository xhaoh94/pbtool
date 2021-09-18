package forms

import (
	"encoding/json"
	"path/filepath"
	"pbtool/conf"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func (f *TMainForm) pbjsPanel(tab *vcl.TTabSheet) {
	cfgs := conf.GetCfgs()
	cfg := cfgs[tab.Tag()]
	pbjsCfg := &conf.PbJsCfg{}
	err := json.Unmarshal([]byte(cfg.Context), pbjsCfg)
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
				updCfg(cfg, pbjsCfg)
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
				updCfg(cfg, pbjsCfg)
			}
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出文件名:")
	lb.SetTop(105)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(p)
	fileEdit.SetParent(p)
	fileEdit.SetBounds(100, 100, 200, 50)
	fileEdit.SetText(pbjsCfg.FileName)
	fileEdit.SetOnExit(func(sender vcl.IObject) {
		if pbjsCfg.FileName != fileEdit.Text() {
			pbjsCfg.FileName = fileEdit.Text()
			updCfg(cfg, pbjsCfg)
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("指定目标格式:")
	lb.SetTop(135)
	lb.SetLeft(10)

	combox1 := vcl.NewComboBox(p)
	combox1.SetParent(p)
	combox1.SetBounds(100, 130, 200, 50)
	combox1.SetStyle(types.CsDropDownList)
	var index int
	for i, v := range conf.PBJS_Target_Comboxs {
		combox1.Items().Add(v)
		if v == pbjsCfg.Target {
			index = i
		}
	}
	combox1.SetItemIndex(int32(index))
	combox1.SetOnChange(func(sender vcl.IObject) {
		target := combox1.Items().Strings(combox1.ItemIndex())
		if pbjsCfg.Target != target {
			pbjsCfg.Target = target
			updCfg(cfg, pbjsCfg)
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("指定包装器:")
	lb.SetTop(165)
	lb.SetLeft(10)

	combox2 := vcl.NewComboBox(p)
	combox2.SetParent(p)
	combox2.SetBounds(100, 160, 200, 50)
	combox2.SetStyle(types.CsDropDownList)
	for i, v := range conf.PBJS_Wrap_Comboxs {
		combox2.Items().Add(v)
		if v == pbjsCfg.Wrap {
			index = i
		}
	}
	combox2.SetItemIndex(int32(index))
	combox2.SetOnChange(func(sender vcl.IObject) {
		wrap := combox2.Items().Strings(combox2.ItemIndex())
		if pbjsCfg.Wrap != wrap {
			pbjsCfg.Wrap = wrap
			updCfg(cfg, pbjsCfg)
		}
	})

	cb := vcl.NewCheckBox(p)
	cb.SetParent(p)
	cb.SetCaption("是否生成d.ts文件")
	cb.SetTop(195)
	cb.SetLeft(10)
	cb.SetChecked(pbjsCfg.CreateDts)
	cb.SetOnChange(func(sender vcl.IObject) {
		if pbjsCfg.CreateDts != cb.Checked() {
			pbjsCfg.CreateDts = cb.Checked()
			updCfg(cfg, pbjsCfg)
		}
	})
	cb2 := vcl.NewCheckBox(p)
	cb2.SetParent(p)
	cb2.SetCaption("是否启用ES6语法")
	cb2.SetTop(225)
	cb2.SetLeft(10)
	cb2.SetChecked(pbjsCfg.UseEs6)
	cb2.SetOnChange(func(sender vcl.IObject) {
		if pbjsCfg.UseEs6 != cb2.Checked() {
			pbjsCfg.UseEs6 = cb2.Checked()
			updCfg(cfg, pbjsCfg)
		}
	})

}

func (f *TCfgForm) pbjsPanel() *PnlCfg {
	pnlCfg := &PnlCfg{
		Cfg: &conf.PbJsCfg{
			Target:    "static-module",
			Wrap:      "commonjs",
			FileName:  "proto.js",
			CreateDts: true,
			UseEs6:    false,
		},
		Pnl: vcl.NewFrame(f.pContext),
	}

	// pnlCfg.Pnl.SetParent(f.pContext)
	pnlCfg.Pnl.SetAlign(types.AlClient)

	lb := vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("导出文件名:")
	lb.SetTop(15)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(pnlCfg.Pnl)
	fileEdit.SetParent(pnlCfg.Pnl)
	fileEdit.SetBounds(100, 10, 200, 50)
	fileEdit.SetText(pnlCfg.Cfg.(*conf.PbJsCfg).FileName)
	fileEdit.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.PbJsCfg).FileName = fileEdit.Text()
	})
	lb = vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("指定目标格式:")
	lb.SetTop(45)
	lb.SetLeft(10)

	combox1 := vcl.NewComboBox(pnlCfg.Pnl)
	combox1.SetParent(pnlCfg.Pnl)
	combox1.SetBounds(100, 40, 200, 50)
	combox1.SetStyle(types.CsDropDownList)
	combox1.SetOnChange(func(sender vcl.IObject) {
		target := combox1.Items().Strings(combox1.ItemIndex())
		pnlCfg.Cfg.(*conf.PbJsCfg).Target = target
	})
	for _, v := range conf.PBJS_Target_Comboxs {
		combox1.Items().Add(v)
	}
	combox1.SetItemIndex(0)

	lb = vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("指定包装器:")
	lb.SetTop(75)
	lb.SetLeft(10)

	combox2 := vcl.NewComboBox(pnlCfg.Pnl)
	combox2.SetParent(pnlCfg.Pnl)
	combox2.SetBounds(100, 70, 200, 50)
	combox2.SetStyle(types.CsDropDownList)
	combox2.SetOnChange(func(sender vcl.IObject) {
		wrap := combox2.Items().Strings(combox2.ItemIndex())
		pnlCfg.Cfg.(*conf.PbJsCfg).Wrap = wrap
	})
	for _, v := range conf.PBJS_Wrap_Comboxs {
		combox2.Items().Add(v)
	}
	combox2.SetItemIndex(0)

	cb := vcl.NewCheckBox(pnlCfg.Pnl)
	cb.SetParent(pnlCfg.Pnl)
	cb.SetCaption("是否生成d.ts文件")
	cb.SetTop(105)
	cb.SetLeft(10)
	cb.SetChecked(pnlCfg.Cfg.(*conf.PbJsCfg).CreateDts)
	cb.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.PbJsCfg).CreateDts = cb.Checked()
	})

	cb2 := vcl.NewCheckBox(pnlCfg.Pnl)
	cb2.SetParent(pnlCfg.Pnl)
	cb2.SetCaption("是否启用ES6语法")
	cb2.SetTop(135)
	cb2.SetLeft(10)
	cb2.SetChecked(pnlCfg.Cfg.(*conf.PbJsCfg).UseEs6)
	cb2.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.PbJsCfg).UseEs6 = cb2.Checked()
	})
	return pnlCfg
}
