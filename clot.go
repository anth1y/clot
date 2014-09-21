package main

import (
	"github.com/bradylove/clot/totp"
	"github.com/nsf/termbox-go"
	"os/exec"
	"strconv"
	"time"
)

const (
	IdColumnWidth     int = 5
	SecretColumnWidth int = 22
	OtpColumnWidth    int = 11

	Password string = "password"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	ts := NewTokenStore("/Users/brady/.clot", Password)

	initView := TokenTable{
		IsVisible:   true,
		SelectedRow: 0,
		TokenStore:  ts,
	}

	initView.ActivateView()
}

type View interface {
	ActivateView()
	DrawView()
}

type TokenTable struct {
	IsVisible   bool
	SelectedRow int
	Tokens      []totp.Token
	TokenStore  TokenStore
}

func (tt *TokenTable) ActivateView() {
	go tt.RefreshLoop()
loop:
	for {
		ev := termbox.PollEvent()

		switch ev.Ch {
		case 'c':
			tt.CopySelectedTokenToClipboard()
		case 'a':
			tt.IsVisible = false
			newToken := NewTokenForm()

			// tt.TokenStore.V1Tokens = append(tt.Tokens, newToken)
			tt.TokenStore.AddToken(newToken)

			tt.DrawView()
		}

		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowDown:
				if tt.SelectedRow < (tt.TokenStore.TokenCount() - 1) {
					tt.SelectedRow++
					tt.DrawView()
				}
			case termbox.KeyArrowUp:
				if tt.SelectedRow > 0 {
					tt.SelectedRow--
					tt.DrawView()
				}
			}
		case termbox.EventResize:
			tt.DrawView()
		}
	}
}

func (tt *TokenTable) RefreshLoop() {
	if tt.IsVisible {
		tt.DrawView()
		time.Sleep(time.Second * 5)
		tt.RefreshLoop()
	}
}

func (tt *TokenTable) DrawOptions() {
	_, termHeight := termbox.Size()

	yPos := termHeight - 2

	xPos := DrawLabel(1, yPos, "c - Copy OTP to clipboard")
	xPos = DrawLabel(xPos+3, yPos, "a - Add token")
	xPos = DrawLabel(xPos+3, yPos, "Esc - Exit")
}

func (tt *TokenTable) CopySelectedTokenToClipboard() {
	token := tt.TokenStore.Tokens()[tt.SelectedRow]

	WriteToClipboard(token.Now())
}

func WriteToClipboard(text string) error {
	copyCmd := exec.Command("reattach-to-user-namespace", "pbcopy")
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}

	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}

	if err := in.Close(); err != nil {
		return err
	}

	return copyCmd.Wait()
}

func (tt *TokenTable) DrawView() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Header & Table Columns
	xPos := DrawColumn(0, IdColumnWidth, "ID", ColumnLeftConnector, ColumnCenterConnector)
	xPos = DrawColumn(xPos, LabelColumnWidth(), "Label", ColumnCenterConnector, ColumnCenterConnector)
	xPos = DrawColumn(xPos, SecretColumnWidth, "Secret", ColumnCenterConnector, ColumnCenterConnector)
	xPos = DrawColumn(xPos, OtpColumnWidth, "OTP", ColumnCenterConnector, ColumnRightConnector)

	// Token Rows
	for index, token := range tt.TokenStore.Tokens() {
		isSelectdToken := index == tt.SelectedRow
		DrawRow(index+1, token, isSelectdToken)
	}

	// Footer
	// tt.DrawFooter()
	tt.DrawOptions()

	termbox.Flush()
	tt.IsVisible = true
}

func (tt *TokenTable) DrawFooter() {
	termWidth, termHeight := termbox.Size()

	yPos := termHeight - 3

	for i := 0; i < termWidth; i++ {
		DrawText(i, yPos, string(HorizontalBar), termbox.ColorDefault, termbox.ColorDefault)
	}
}

func DrawColumn(xPos, width int, headText string, startChar, endChar Symbol) int {
	_, termHeight := termbox.Size()

	termHeight = termHeight - 3

	DrawText(xPos, 0, string(VerticalBar), termbox.ColorBlack, termbox.ColorWhite)
	for i := 1; i < termHeight; i++ {
		DrawLabel(xPos, i, string(VerticalBar))
	}
	DrawLabel(xPos, termHeight, string(startChar))

	xPos++

	oldXPos := xPos

	xPos = DrawHeader(xPos, 0, RightPadString(width, " "+headText))
	DrawText(xPos, 0, string(VerticalBar), termbox.ColorBlack, termbox.ColorWhite)
	for i := 1; i < termHeight; i++ {
		DrawLabel(xPos, i, string(VerticalBar))
	}

	DrawLabel(xPos, termHeight, string(endChar))
	for i := 0; i < (xPos - oldXPos); i++ {
		DrawLabel(i+oldXPos, termHeight, string(HorizontalBar))
	}

	return xPos
}

func RightPadString(size int, text string) string {
	padSize := size - len(text)

	for i := 0; i < padSize; i++ {
		text = text + " "
	}

	return text
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

	xPos = DrawText(xPos, yPos, string(VerticalBar), fg, bg)
	xPos = DrawText(xPos, yPos, RightPadString(IdColumnWidth, " "+strconv.Itoa(t.Id)), fg, bg)
	xPos = DrawText(xPos, yPos, string(VerticalBar), fg, bg)
	xPos = DrawText(xPos, yPos, RightPadString(LabelColumnWidth(), " "+t.Label), fg, bg)
	xPos = DrawText(xPos, yPos, string(VerticalBar), fg, bg)
	xPos = DrawText(xPos, yPos, RightPadString(SecretColumnWidth, " "+t.Secret), fg, bg)
	xPos = DrawText(xPos, yPos, string(VerticalBar), fg, bg)
	xPos = DrawText(xPos, yPos, RightPadString(OtpColumnWidth, " "+t.Now()), fg, bg)
	xPos = DrawText(xPos, yPos, string(VerticalBar), fg, bg)
}

func LabelColumnWidth() int {
	termWidth, _ := termbox.Size()
	labelWidth := termWidth - (IdColumnWidth + SecretColumnWidth + OtpColumnWidth + 5)

	return labelWidth
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
