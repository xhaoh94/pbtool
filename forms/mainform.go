package forms

import (
	"encoding/json"
	"path/filepath"
	"pbtool/conf"
	"pbtool/parse"
	"sort"

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
		b := parse.Parse(cfg)
		if b {
			vcl.ShowMessage("导出成功")
		}
	})

	btn2 := vcl.NewButton(pBottom)
	btn2.SetParent(pBottom)
	btn2.SetCaption("导出所有配置")
	btn2.SetBounds(220, 25, 200, 50)
	btn2.SetOnClick(func(sender vcl.IObject) {
		cfgs := conf.GetCfgs()
		var b bool
		for key := range cfgs {
			cfg := cfgs[key]
			b = parse.Parse(cfg)
			if !b {
				break
			}
		}

		if b {
			vcl.ShowMessage("导出成功")
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
	var keys []int
	for k := range cfgs {
		keys = append(keys, cfgs[k].ID)
	}
	sort.Ints(keys)
	for _, k := range keys {
		v := cfgs[k]
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
			f.tsPanel(tab)
			break
		case conf.TC_Golang:
			f.goPanel(tab)
			break
		case conf.TC_CSharp:
			f.cSharpPanel(tab)
			break
		case conf.TC_Protobufjs:
			f.pbjsPanel(tab)
			break
		}
	}
}

func updCfg(cfg *conf.OutCfg, context interface{}) {
	data, _ := json.Marshal(context)
	cfg.Context = string(data)
	conf.AddCfg(cfg)
}

func onBtnAddCfgClick(sender vcl.IObject) {
	CfgForm.Init()
	CfgForm.Show()
}
