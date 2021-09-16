package forms

import (
	"encoding/json"
	"hash/crc32"
	"path/filepath"
	"pbtool/conf"

	"github.com/google/uuid"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
)

type TCfgForm struct {
	*vcl.TForm
	pContext *vcl.TPanel
	cfg      *conf.OutCfg

	ts *conf.TsCfg
}

var (
	CfgForm *TCfgForm
)

func (f *TCfgForm) OnFormCreate(sender vcl.IObject) {
	f.cfg = new(conf.OutCfg)
	f.cfg.ID = StringToHash(uuid.New().String())
	f.SetCaption("添加配置")
	f.SetBorderStyle(types.BsToolWindow)
	f.SetPosition(types.PoScreenCenter)
	f.SetWidth(600)
	f.SetHeight(300)
	f.create()
	f.tsPanel()
}
func (f *TCfgForm) create() {

	pTop := vcl.NewFrame(f)
	pTop.SetParent(f)
	// pTop.SetParentBackground(false)
	// pTop.SetColor(colors.ClWhite)
	pTop.SetHeight(120)
	pTop.SetAlign(types.AlTop)

	lb := vcl.NewLabel(pTop)
	lb.SetParent(pTop)
	lb.SetCaption("proto文件目录：")
	lb.SetTop(15)
	lb.SetLeft(10)

	dirEdit := vcl.NewEdit(pTop)
	dirEdit.SetEnabled(false)
	dirEdit.SetParent(pTop)
	dirEdit.SetBounds(100, 10, 400, 50)

	btnDir := vcl.NewButton(pTop)
	btnDir.SetParent(pTop)
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

	lb = vcl.NewLabel(pTop)
	lb.SetParent(pTop)
	lb.SetCaption("导出文件目录：")
	lb.SetTop(50)
	lb.SetLeft(10)

	outEdit := vcl.NewEdit(pTop)
	outEdit.SetEnabled(false)
	outEdit.SetParent(pTop)
	outEdit.SetBounds(100, 50, 400, 50)

	btnOut := vcl.NewButton(pTop)
	btnOut.SetParent(pTop)
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

	lb = vcl.NewLabel(pTop)
	lb.SetParent(pTop)
	lb.SetCaption("配置名称:")
	lb.SetTop(95)
	lb.SetLeft(10)

	nameEdit := vcl.NewEdit(pTop)
	nameEdit.SetParent(pTop)
	nameEdit.SetTextHint("请输入名字")
	nameEdit.SetBounds(100, 90, 200, 50)
	nameEdit.SetOnChange(func(sender vcl.IObject) {
		f.cfg.Name = nameEdit.Text()
	})

	lb = vcl.NewLabel(pTop)
	lb.SetParent(pTop)
	lb.SetCaption("类型:")
	lb.SetTop(95)
	lb.SetLeft(310)

	typeCombox := vcl.NewComboBox(pTop)
	typeCombox.SetParent(pTop)
	typeCombox.SetTop(90)
	typeCombox.SetLeft(350)
	typeCombox.SetStyle(types.CsDropDownList)
	typeCombox.SetOnChange(func(sender vcl.IObject) {
		f.cfg.TagType = typeCombox.Items().Strings(typeCombox.ItemIndex())
	})
	for _, v := range conf.TypeComboxs {
		typeCombox.Items().Add(v)
	}
	typeCombox.SetItemIndex(0)
	f.cfg.TagType = typeCombox.Items().Strings(typeCombox.ItemIndex())

	btnSave := vcl.NewButton(pTop)
	btnSave.SetParent(pTop)
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
		var err error
		var data []byte
		switch f.cfg.TagType {
		case conf.TC_TypeScript:
			data, err = json.Marshal(f.ts)
			break
		}
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

	f.pContext = vcl.NewPanel(f)
	f.pContext.SetParent(f)
	f.pContext.SetParentBackground(false)
	f.pContext.SetColor(colors.ClWhite)
	f.pContext.SetAlign(types.AlClient)

}

func (f *TCfgForm) tsPanel() {
	f.ts = &conf.TsCfg{
		Ns:          "proto",
		CreateJson:  true,
		UseModule:   true,
		OutJsonPath: "",
		FileName:    "DymPbCode.ts",
		JsonName:    "DymPbCfg.json",
	}
	p := vcl.NewFrame(f.pContext)
	p.SetParent(f.pContext)
	p.SetAlign(types.AlClient)

	lb := vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("命名空间:")
	lb.SetTop(15)
	lb.SetLeft(10)

	edit := vcl.NewEdit(p)
	edit.SetParent(p)
	edit.SetBounds(100, 10, 200, 50)
	edit.SetText(f.ts.Ns)
	edit.SetOnChange(func(sender vcl.IObject) {
		f.ts.Ns = edit.Text()
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出文件名:")
	lb.SetTop(45)
	lb.SetLeft(10)

	fileEdit := vcl.NewEdit(p)
	fileEdit.SetParent(p)
	fileEdit.SetBounds(100, 40, 200, 50)
	fileEdit.SetText(f.ts.FileName)
	fileEdit.SetOnChange(func(sender vcl.IObject) {
		f.ts.FileName = fileEdit.Text()
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出Json名:")
	lb.SetTop(75)
	lb.SetLeft(10)

	jsonNameEdit := vcl.NewEdit(p)
	jsonNameEdit.SetParent(p)
	jsonNameEdit.SetBounds(100, 70, 200, 50)
	jsonNameEdit.SetText(f.ts.JsonName)
	jsonNameEdit.SetOnChange(func(sender vcl.IObject) {
		f.ts.JsonName = jsonNameEdit.Text()
	})

	lb = vcl.NewLabel(p)
	lb.SetParent(p)
	lb.SetCaption("导出Json目录：")
	lb.SetTop(105)
	lb.SetLeft(10)

	jsonEdit := vcl.NewEdit(p)
	jsonEdit.SetEnabled(false)
	jsonEdit.SetParent(p)
	jsonEdit.SetBounds(100, 100, 400, 50)

	btnJson := vcl.NewButton(p)
	btnJson.SetParent(p)
	btnJson.SetCaption("···")
	btnJson.SetBounds(510, 100, 80, 25)
	btnJson.SetOnClick(func(sender vcl.IObject) {
		sdd := vcl.NewSelectDirectoryDialog(nil)
		if sdd.Execute() {
			path := filepath.ToSlash(sdd.FileName())
			jsonEdit.SetText(path)
			f.ts.OutJsonPath = path
		}
	})

	cb1 := vcl.NewCheckBox(p)
	cb1.SetParent(p)
	cb1.SetCaption("是否生成Json配置")
	cb1.SetTop(135)
	cb1.SetLeft(10)
	cb1.SetChecked(f.ts.CreateJson)
	cb1.SetOnChange(func(sender vcl.IObject) {
		f.ts.CreateJson = cb1.Checked()
	})

	cb2 := vcl.NewCheckBox(p)
	cb2.SetParent(p)
	cb2.SetCaption("是否增加export关键字")
	cb2.SetTop(155)
	cb2.SetLeft(10)
	cb2.SetChecked(f.ts.UseModule)
	cb2.SetOnChange(func(sender vcl.IObject) {
		f.ts.UseModule = cb2.Checked()
	})
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
