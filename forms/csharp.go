package forms

import (
	"encoding/json"
	"path/filepath"
	"pbtool/conf"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func (f *TMainForm) cSharpPanel(tab *vcl.TTabSheet) {
	cfgs := conf.GetCfgs()
	cfg := cfgs[tab.Tag()]
	cSharpCfg := &conf.CSharpCfg{}
	err := json.Unmarshal([]byte(cfg.Context), cSharpCfg)
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
				updCfg(cfg, cSharpCfg)
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
				updCfg(cfg, cSharpCfg)
			}
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("命名空间:")
	lb.SetTop(105)
	lb.SetLeft(10)

	edit := vcl.NewEdit(p)
	edit.SetParent(p)
	edit.SetBounds(100, 100, 200, 50)
	edit.SetOnExit(func(sender vcl.IObject) {
		if cSharpCfg.Ns != edit.Text() {
			cSharpCfg.Ns = edit.Text()
			updCfg(cfg, cSharpCfg)
		}
	})
	edit.SetText(cSharpCfg.Ns)

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出文件名:")
	lb.SetTop(135)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(p)
	fileEdit.SetParent(p)
	fileEdit.SetBounds(100, 130, 200, 50)
	fileEdit.SetText(cSharpCfg.FileName)
	fileEdit.SetOnExit(func(sender vcl.IObject) {
		if cSharpCfg.FileName != fileEdit.Text() {
			cSharpCfg.FileName = fileEdit.Text()
			updCfg(cfg, cSharpCfg)
		}
	})

	cb1 := vcl.NewCheckBox(p)
	cb1.SetParent(p)
	cb1.SetCaption("是否生成CMD文件")
	cb1.SetTop(175)
	cb1.SetLeft(10)
	cb1.SetChecked(cSharpCfg.CreateCmd)
	cb1.SetOnChange(func(sender vcl.IObject) {
		if cSharpCfg.CreateCmd != cb1.Checked() {
			cSharpCfg.CreateCmd = cb1.Checked()
			updCfg(cfg, cSharpCfg)
		}

	})

}

func (f *TCfgForm) cSharpPanel() *PnlCfg {
	pnlCfg := &PnlCfg{
		Cfg: &conf.CSharpCfg{
			Ns:        "pb",
			CreateCmd: true,
			FileName:  "Cmd.cs",
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
	edit.SetText(pnlCfg.Cfg.(*conf.CSharpCfg).Ns)
	edit.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.CSharpCfg).Ns = edit.Text()
	})

	lb = vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("导出文件名:")
	lb.SetTop(45)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(pnlCfg.Pnl)
	fileEdit.SetParent(pnlCfg.Pnl)
	fileEdit.SetBounds(100, 40, 200, 50)
	fileEdit.SetText(pnlCfg.Cfg.(*conf.CSharpCfg).FileName)
	fileEdit.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.CSharpCfg).FileName = fileEdit.Text()
	})

	cb1 := vcl.NewCheckBox(pnlCfg.Pnl)
	cb1.SetParent(pnlCfg.Pnl)
	cb1.SetCaption("是否生成CMD文件")
	cb1.SetTop(75)
	cb1.SetLeft(10)
	cb1.SetChecked(pnlCfg.Cfg.(*conf.CSharpCfg).CreateCmd)
	cb1.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.CSharpCfg).CreateCmd = cb1.Checked()
	})
	return pnlCfg
}
