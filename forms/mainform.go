package forms

import (
	"encoding/json"
	"path/filepath"
	"pbtool/conf"
	"pbtool/parse"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type TMainForm struct {
	*vcl.TForm

	pageCol *vcl.TPageControl
}

var (
	MainForm *TMainForm

	id2tab map[int]*vcl.TTabSheet = make(map[int]*vcl.TTabSheet)
)

func (f *TMainForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("protobufTool")
	f.SetBorderStyle(types.BsSingle)
	f.SetPosition(types.PoScreenCenter)
	f.SetWidth(720)
	f.SetHeight(500)
	f.createTop()
	f.createPage()
	f.createBottom()
}

func (f *TMainForm) createTop() {
	pTop := vcl.NewFrame(f)
	pTop.SetParent(f)
	pTop.SetHeight(100)
	// pTop.SetColor(colors.ClYellow)
	pTop.SetAlign(types.AlTop)
	lb := vcl.NewLabel(pTop)
	lb.SetParent(pTop)
	lb.SetCaption("protoc.exe路径：")
	lb.SetTop(15)
	lb.SetLeft(10)

	edit := vcl.NewEdit(pTop)
	edit.SetEnabled(false)
	edit.SetParent(pTop)
	edit.SetBounds(120, 10, 480, 50)
	edit.SetText(conf.GetProtoPath())

	btnOut := vcl.NewButton(pTop)
	btnOut.SetParent(pTop)
	btnOut.SetCaption("···")
	btnOut.SetBounds(620, 10, 80, 25)
	btnOut.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewOpenDialog(nil)
		sdd.SetFilter("protoc(*.exe)|*.exe")
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			edit.SetText(path)
			conf.SetProtoPath(path)
		}
	})

	btnAdd := vcl.NewButton(pTop)
	btnAdd.SetParent(pTop)
	btnAdd.SetCaption("添加配置")
	btnAdd.SetBounds(10, 60, 100, 30)
	btnAdd.SetOnClick(onBtnAddCfgClick)

	btnDel := vcl.NewButton(pTop)
	btnDel.SetParent(pTop)
	btnDel.SetCaption("删除配置")
	btnDel.SetBounds(120, 60, 100, 30)
	btnDel.SetOnClick(func(sen vcl.IObject) {
		index := f.pageCol.ActivePageIndex()
		tab := f.pageCol.Controls(index)
		id := tab.Tag()
		conf.DelCfg(id)
		f.updPage()
	})
}

func (f *TMainForm) createBottom() {
	pBottom := vcl.NewPanel(f)
	pBottom.SetParent(f)
	pBottom.SetHeight(100)
	pBottom.SetAlign(types.AlBottom)

	btn1 := vcl.NewButton(pBottom)
	btn1.SetParent(pBottom)
	btn1.SetCaption("导出当前配置")
	btn1.SetBounds(10, 25, 200, 50)
	btn1.SetOnClick(func(sender vcl.IObject) {
		index := f.pageCol.ActivePageIndex()
		tab := f.pageCol.Controls(index)
		id := tab.Tag()
		cfg, ok := conf.GetCfgs()[id]
		if !ok {
			return
		}
		parse.Parse(cfg)
	})

	btn2 := vcl.NewButton(pBottom)
	btn2.SetParent(pBottom)
	btn2.SetCaption("导出所有配置")
	btn2.SetBounds(220, 25, 200, 50)
	btn2.SetOnClick(func(sender vcl.IObject) {
		cfgs := conf.GetCfgs()
		for key := range cfgs {
			cfg := cfgs[key]
			parse.Parse(cfg)
		}
	})
}

func (f *TMainForm) createPage() {
	f.pageCol = vcl.NewPageControl(f)
	f.pageCol.SetParent(f)
	// 这里将TPageControl设置为整个窗口客户区大小，并自动调整
	f.pageCol.SetAlign(types.AlClient)

	f.updPage()
}

func (f *TMainForm) updPage() {
	cfgs := conf.GetCfgs()
	if len(id2tab) > 0 {
		for k := range id2tab {
			v := id2tab[k]
			if cfgs[k] == nil {
				delete(id2tab, k)
				v.SetPageControl(nil)
			}
		}
	}
	for _, v := range cfgs {
		tab := id2tab[v.ID]
		if tab == nil {
			tab = vcl.NewTabSheet(f.pageCol)
			tab.SetTag(v.ID)
			tab.SetPageControl(f.pageCol)
			tab.SetCaption(v.Name)
			id2tab[v.ID] = tab
		}
		switch v.TagType {
		case conf.TC_TypeScript:
			tsPanel(tab)
			break
		}
	}
}

func tsPanel(tab *vcl.TTabSheet) {
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

func updCfg(cfg *conf.OutCfg, context interface{}) {
	data, _ := json.Marshal(context)
	cfg.Context = string(data)
	conf.AddCfg(cfg)
}

func onBtnAddCfgClick(sender vcl.IObject) {
	CfgForm.Show()
}
