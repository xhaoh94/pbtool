package forms

import (
	"encoding/json"
	"hash/crc32"
	"path/filepath"
	"pbtool/conf"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
)

type TCfgForm struct {
	*vcl.TForm
	pTop     *vcl.TFrame
	pContext *vcl.TPanel
	cfg      *conf.OutCfg

	tag2pnl map[string]*PnlCfg
}
type PnlCfg struct {
	Pnl *vcl.TFrame
	Cfg interface{}
}

var (
	CfgForm *TCfgForm
)

func (f *TCfgForm) Init() {
	f.cfg = new(conf.OutCfg)
	l := len(conf.GetCfgs())
	f.cfg.ID = l + 1
	f.tag2pnl = make(map[string]*PnlCfg)
	f.create()
}

func (f *TCfgForm) OnFormCreate(sender vcl.IObject) {
	f.SetCaption("添加配置")
	f.SetBorderStyle(types.BsSingle)
	f.SetBorderIcons(1)
	f.SetWidth(600)
	f.SetHeight(300)
	f.SetPosition(types.PoScreenCenter)
}
func (f *TCfgForm) create() {
	if f.pContext != nil {
		f.pContext.SetParent(nil)
		f.pContext.Free()
		f.pContext = nil
	}
	if f.pTop != nil {
		f.pTop.SetParent(nil)
		f.pTop.Free()
		f.pTop = nil
	}

	f.pContext = vcl.NewPanel(f)
	f.pContext.SetParent(f)
	f.pContext.SetParentBackground(false)
	f.pContext.SetColor(colors.ClWhite)
	f.pContext.SetAlign(types.AlClient)

	f.pTop = vcl.NewFrame(f)
	f.pTop.SetParent(f)
	// f.pTop.SetParentBackground(false)
	// f.pTop.SetColor(colors.ClWhite)
	f.pTop.SetHeight(120)
	f.pTop.SetAlign(types.AlTop)

	lb := vcl.NewLabel(f.pTop)
	lb.SetParent(f.pTop)
	lb.SetCaption("proto文件目录：")
	lb.SetTop(15)
	lb.SetLeft(10)

	dirEdit := vcl.NewEdit(f.pTop)
	dirEdit.SetEnabled(false)
	dirEdit.SetParent(f.pTop)
	dirEdit.SetBounds(100, 10, 400, 50)

	btnDir := vcl.NewButton(f.pTop)
	btnDir.SetParent(f.pTop)
	btnDir.SetCaption("···")
	btnDir.SetBounds(510, 10, 80, 25)
	btnDir.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewSelectDirectoryDialog(nil)
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			dirEdit.SetText(path)
			f.cfg.InPath = dirEdit.Text()
		}
	})

	lb = vcl.NewLabel(f.pTop)
	lb.SetParent(f.pTop)
	lb.SetCaption("导出文件目录：")
	lb.SetTop(50)
	lb.SetLeft(10)

	outEdit := vcl.NewEdit(f.pTop)
	outEdit.SetEnabled(false)
	outEdit.SetParent(f.pTop)
	outEdit.SetBounds(100, 50, 400, 50)

	btnOut := vcl.NewButton(f.pTop)
	btnOut.SetParent(f.pTop)
	btnOut.SetCaption("···")
	btnOut.SetBounds(510, 50, 80, 25)
	btnOut.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewSelectDirectoryDialog(nil)
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			outEdit.SetText(path)
			f.cfg.OutPath = path
		}
	})

	lb = vcl.NewLabel(f.pTop)
	lb.SetParent(f.pTop)
	lb.SetCaption("配置名称:")
	lb.SetTop(95)
	lb.SetLeft(10)

	nameEdit := vcl.NewEdit(f.pTop)
	nameEdit.SetParent(f.pTop)
	nameEdit.SetTextHint("请输入名字")
	nameEdit.SetBounds(100, 90, 200, 50)
	nameEdit.SetOnChange(func(sender vcl.IObject) {
		f.cfg.Name = nameEdit.Text()
	})

	lb = vcl.NewLabel(f.pTop)
	lb.SetParent(f.pTop)
	lb.SetCaption("类型:")
	lb.SetTop(95)
	lb.SetLeft(310)

	typeCombox := vcl.NewComboBox(f.pTop)
	typeCombox.SetParent(f.pTop)
	typeCombox.SetTop(90)
	typeCombox.SetLeft(350)
	typeCombox.SetStyle(types.CsDropDownList)
	typeCombox.SetOnChange(func(sender vcl.IObject) {
		TagType := typeCombox.Items().Strings(typeCombox.ItemIndex())
		f.comboxChanged(TagType)
	})
	for _, v := range conf.TypeComboxs {
		typeCombox.Items().Add(v)
		switch v {
		case conf.TC_TypeScript:
			f.tag2pnl[v] = f.tsPanel()
			break
		case conf.TC_Golang:
			f.tag2pnl[v] = f.goPanel()
			break
		case conf.TC_CSharp:
			f.tag2pnl[v] = f.cSharpPanel()
			break
		case conf.TC_Protobufjs:
			f.tag2pnl[v] = f.pbjsPanel()
			break
		}
	}
	typeCombox.SetItemIndex(0)
	f.comboxChanged(typeCombox.Items().Strings(typeCombox.ItemIndex()))

	btnSave := vcl.NewButton(f.pTop)
	btnSave.SetParent(f.pTop)
	btnSave.SetCaption("保存")
	btnSave.SetBounds(510, 90, 80, 25)
	btnSave.SetOnClick(func(sender vcl.IObject) {
		if f.cfg.Name == "" {
			vcl.ShowMessage("配置名不能为空")
			return
		}
		if f.cfg.InPath == "" {
			vcl.ShowMessage("proto文件路径不能为空")
			return
		}
		if f.cfg.OutPath == "" {
			vcl.ShowMessage("导出路径不能为空")
			return
		}
		pnlCfg := f.tag2pnl[f.cfg.TagType]
		data, err := json.Marshal(pnlCfg.Cfg)
		if err != nil {
			vcl.ShowMessage("解析Json错误")
			return
		}

		f.cfg.Context = string(data)
		conf.AddCfg(f.cfg)
		MainForm.updPage()
		f.Close()
		MainForm.SetFocus()
	})

}

func (f *TCfgForm) comboxChanged(TagType string) {
	if f.cfg.TagType == TagType {
		return
	}
	if f.cfg.TagType != "" {
		pnlCfg := f.tag2pnl[f.cfg.TagType]
		// pnlCfg.Pnl.SetVisible(false)
		pnlCfg.Pnl.SetParent(nil)
	}
	f.cfg.TagType = TagType
	pnlCfg := f.tag2pnl[TagType]
	// pnlCfg.Pnl.SetVisible(true)
	pnlCfg.Pnl.SetParent(f.pContext)
}

//StringToHash 字符串转为32位整形哈希
func StringToHash(s string) (hash int) {

	hash = int(crc32.ChecksumIEEE([]byte(s)))
	if hash >= 0 {
		return hash
	}
	if -hash >= 0 {
		return -hash
	}

	for _, c := range s {
		ch := int(c)
		hash = hash + ((hash) << 5) + ch + (ch << 7)
	}
	return
}
