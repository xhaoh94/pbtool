package pbjs

type (
	MessageStruct struct {
		Cs    uint32
		Sc    uint32
		Title string
		Datas [][]string
	}
)

var (
	Messages []*MessageStruct
)
