package main

import (
	"clot/totp"
	"github.com/nsf/termbox-go"
	"strconv"
	"time"
)

const (
	IdColumnWidth     int = 5
	SecretColumnWidth int = 22
	OtpColumnWidth    int = 11
)

var tokens []totp.Token
var selectedToken int

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	tokens = make([]totp.Token, 0)
	tokens = append(tokens, totp.Token{
		Id:       1,
		Label:    "Brady.Love@nccgroup.com",
		Secret:   "ptcuyfm2sjjqgh5c",
		Digits:   6,
		TimeStep: 30,
	})
	tokens = append(tokens, totp.Token{
		Id:       2,
		Label:    "Brady.Love@gmail.com",
		Secret:   "4wd3tgngs65ybcps",
		Digits:   6,
		TimeStep: 30,
	})
	tokens = append(tokens, totp.Token{
		Id:       3,
		Label:    "all token all the time",
		Secret:   "zrbxo6diith5wqwd",
		Digits:   6,
		TimeStep: 30,
	})
	tokens = append(tokens, totp.Token{
		Id:       4,
		Label:    "tokin it token",
		Secret:   "r3gtagvpj76q5j73",
		Digits:   6,
		TimeStep: 30,
	})
	tokens = append(tokens, totp.Token{
		Id:       5,
		Label:    "much token for much otp",
		Secret:   "4wd3tgngs65ybcps",
		Digits:   6,
		TimeStep: 30,
	})
	tokens = append(tokens, totp.Token{
		Id:       6,
		Label:    "just another token",
		Secret:   "mz2cypmurwfie5kx",
		Digits:   6,
		TimeStep: 30,
	})

	selectedToken = 0

	go refresher()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowDown:
				if selectedToken < (len(tokens) - 1) {
					selectedToken++
					DrawTokenTable(selectedToken)
				}
			case termbox.KeyArrowUp:
				if selectedToken > 0 {
					selectedToken--
					DrawTokenTable(selectedToken)
				}
			case termbox.KeyEsc:
				break loop
			}
		case termbox.EventResize:
			DrawTokenTable(0)
		}
	}
}

func refresher() {
	DrawTokenTable(selectedToken)
	time.Sleep(time.Second * 5)

	refresher()
}

func DrawNewTokenForm() {

}

func DrawTokenTable(selectedRow int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	xPos := DrawColumn(0, IdColumnWidth, "ID")
	xPos = DrawColumn(xPos, LabelColumnWidth(), "Label")
	xPos = DrawColumn(xPos, SecretColumnWidth, "Secret")
	xPos = DrawColumn(xPos, OtpColumnWidth, "OTP")

	for index, token := range tokens {
		isSelectdToken := index == selectedRow
		DrawRow(index+1, token, isSelectdToken)
	}

	termbox.Flush()
}

func DrawColumn(xPos, width int, headText string) int {
	_, termHeight := termbox.Size()

	DrawText(xPos, 0, string(0x2502), termbox.ColorBlack, termbox.ColorWhite)
	for i := 1; i < termHeight; i++ {
		DrawLabel(xPos, i, string(0x2502))
	}

	xPos++

	xPos = DrawHeader(xPos, 0, RightPadString(width, " "+headText))
	DrawText(xPos, 0, string(0x2502), termbox.ColorBlack, termbox.ColorWhite)
	for i := 1; i < termHeight; i++ {
		DrawLabel(xPos, i, string(0x2502))
	}

	return xPos
}

func DrawRow(yPos int, t totp.Token, selected bool) {
	xPos := 0

	var fg, bg termbox.Attribute

	if selected {
		bg = termbox.ColorGreen
		fg = termbox.ColorBlack
	} else {
		bg = termbox.ColorDefault
		fg = termbox.ColorDefault
	}

	xPos = DrawText(xPos, yPos, string(0x2502), fg, bg)
	xPos = DrawText(xPos, yPos, RightPadString(IdColumnWidth, " "+strconv.Itoa(t.Id)), fg, bg)
	xPos = DrawText(xPos, yPos, string(0x2502), fg, bg)
	xPos = DrawText(xPos, yPos, RightPadString(LabelColumnWidth(), " "+t.Label), fg, bg)
	xPos = DrawText(xPos, yPos, string(0x2502), fg, bg)
	xPos = DrawText(xPos, yPos, RightPadString(SecretColumnWidth, " "+t.Secret), fg, bg)
	xPos = DrawText(xPos, yPos, string(0x2502), fg, bg)
	xPos = DrawText(xPos, yPos, RightPadString(OtpColumnWidth, " "+t.Now()), fg, bg)
	xPos = DrawText(xPos, yPos, string(0x2502), fg, bg)
}

func LabelColumnWidth() int {
	termWidth, _ := termbox.Size()
	labelWidth := termWidth - 43

	return labelWidth
}

func RightPadString(size int, text string) string {
	padSize := size - len(text)

	for i := 0; i < padSize; i++ {
		text = text + " "
	}

	return text
}

func DrawHeader(x, y int, text string) int {
	return DrawText(x, y, text, termbox.ColorBlack, termbox.ColorWhite)
}

func DrawLabel(x, y int, text string) int {
	return DrawText(x, y, text, termbox.ColorDefault, termbox.ColorDefault)
}

func DrawText(x, y int, text string, fg, bg termbox.Attribute) int {
	for _, rv := range text {
		termbox.SetCell(x, y, rv, fg, bg)
		x++
	}

	return x
}
