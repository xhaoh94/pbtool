package forms

import (
	"encoding/json"
	"path/filepath"
	"pbtool/conf"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func (f *TMainForm) goPanel(tab *vcl.TTabSheet) {
	cfgs := conf.GetCfgs()
	cfg := cfgs[tab.Tag()]
	goCfg := &conf.GoCfg{}
	err := json.Unmarshal([]byte(cfg.Context), goCfg)
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
				updCfg(cfg, goCfg)
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
				updCfg(cfg, goCfg)
			}
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("proto-gen-go.exe路径：")
	lb.SetTop(105)
	lb.SetLeft(10)

	protoGengoEdit := vcl.NewEdit(p)
	protoGengoEdit.SetEnabled(false)
	protoGengoEdit.SetParent(p)
	protoGengoEdit.SetBounds(10, 130, 590, 50)
	protoGengoEdit.SetText(goCfg.ProtoGenGoPath)

	btnProtoGengo := vcl.NewButton(p)
	btnProtoGengo.SetParent(p)
	btnProtoGengo.SetCaption("···")
	btnProtoGengo.SetBounds(620, 130, 80, 25)
	btnProtoGengo.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewOpenDialog(nil)
		sdd.SetFilter("proto-gen-go(*.exe)|*.exe")
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			if goCfg.ProtoGenGoPath != path {
				protoGengoEdit.SetText(path)
				goCfg.ProtoGenGoPath = path
				updCfg(cfg, goCfg)
			}
		}
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("命名空间:")
	lb.SetTop(165)
	lb.SetLeft(10)

	edit := vcl.NewEdit(p)
	edit.SetParent(p)
	edit.SetBounds(100, 160, 200, 50)
	edit.SetOnExit(func(sender vcl.IObject) {
		if goCfg.Ns != edit.Text() {
			goCfg.Ns = edit.Text()
			updCfg(cfg, goCfg)
		}
	})
	edit.SetText(goCfg.Ns)

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出文件名:")
	lb.SetTop(195)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(p)
	fileEdit.SetParent(p)
	fileEdit.SetBounds(100, 190, 200, 50)
	fileEdit.SetText(goCfg.FileName)
	fileEdit.SetOnExit(func(sender vcl.IObject) {
		if goCfg.FileName != fileEdit.Text() {
			goCfg.FileName = fileEdit.Text()
			updCfg(cfg, goCfg)
		}
	})

	cb1 := vcl.NewCheckBox(p)
	cb1.SetParent(p)
	cb1.SetCaption("是否生成CMD文件")
	cb1.SetTop(225)
	cb1.SetLeft(10)
	cb1.SetChecked(goCfg.CreateCmd)
	cb1.SetOnChange(func(sender vcl.IObject) {
		if goCfg.CreateCmd != cb1.Checked() {
			goCfg.CreateCmd = cb1.Checked()
			updCfg(cfg, goCfg)
		}

	})

}

func (f *TCfgForm) goPanel() *PnlCfg {
	pnlCfg := &PnlCfg{
		Cfg: &conf.GoCfg{
			Ns:             "pb",
			CreateCmd:      true,
			FileName:       "cmd.pb.go",
			ProtoGenGoPath: "",
		},
		Pnl: vcl.NewFrame(f.pContext),
	}

	// pnlCfg.Pnl.SetParent(f.pContext)
	pnlCfg.Pnl.SetAlign(types.AlClient)

	lb := vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("protoc-gen-go.exe路径：")
	lb.SetTop(15)
	lb.SetLeft(10)

	protoGenPathEdit := vcl.NewEdit(pnlCfg.Pnl)
	protoGenPathEdit.SetEnabled(false)
	protoGenPathEdit.SetParent(pnlCfg.Pnl)
	protoGenPathEdit.SetBounds(10, 40, 490, 50)

	btnProtoGenPath := vcl.NewButton(pnlCfg.Pnl)
	btnProtoGenPath.SetParent(pnlCfg.Pnl)
	btnProtoGenPath.SetCaption("···")
	btnProtoGenPath.SetBounds(510, 40, 80, 25)
	btnProtoGenPath.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewOpenDialog(nil)
		sdd.SetFilter("proto-gen-go(*.exe)|*.exe")
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			protoGenPathEdit.SetText(path)
			pnlCfg.Cfg.(*conf.GoCfg).ProtoGenGoPath = path
		}
	})

	lb = vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("命名空间:")
	lb.SetTop(75)
	lb.SetLeft(10)

	edit := vcl.NewEdit(pnlCfg.Pnl)
	edit.SetParent(pnlCfg.Pnl)
	edit.SetBounds(100, 70, 200, 50)
	edit.SetText(pnlCfg.Cfg.(*conf.GoCfg).Ns)
	edit.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.GoCfg).Ns = edit.Text()
	})

	lb = vcl.NewLabel(pnlCfg.Pnl)
	lb.SetParent(pnlCfg.Pnl)
	lb.SetCaption("导出文件名:")
	lb.SetTop(105)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(pnlCfg.Pnl)
	fileEdit.SetParent(pnlCfg.Pnl)
	fileEdit.SetBounds(100, 100, 200, 50)
	fileEdit.SetText(pnlCfg.Cfg.(*conf.GoCfg).FileName)
	fileEdit.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.GoCfg).FileName = fileEdit.Text()
	})

	cb1 := vcl.NewCheckBox(pnlCfg.Pnl)
	cb1.SetParent(pnlCfg.Pnl)
	cb1.SetCaption("是否生成CMD文件")
	cb1.SetTop(135)
	cb1.SetLeft(10)
	cb1.SetChecked(pnlCfg.Cfg.(*conf.GoCfg).CreateCmd)
	cb1.SetOnChange(func(sender vcl.IObject) {
		pnlCfg.Cfg.(*conf.GoCfg).CreateCmd = cb1.Checked()
	})
	return pnlCfg
}
