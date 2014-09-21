package main

import (
	termadd "github.com/Chownie/Termbox-Additions"
	termaddutils "github.com/Chownie/Termbox-Additions/utils"
	"github.com/bradylove/clot/totp"
	"github.com/nsf/termbox-go"
)

func NewTokenForm() totp.Token {
	termWidth, termHeight := termbox.Size()

	xPos := (termWidth / 2) - 50
	yPos := (termHeight / 2) - 2

	label := termadd.DrawForm(xPos, yPos, "Token Label", termadd.AL_LEFT, termaddutils.CONNECT_NONE, 100)
	secret := termadd.DrawForm(xPos, yPos, "Token Secret", termadd.AL_LEFT, termaddutils.CONNECT_NONE, 100)

	return totp.Token{
		Id:       1,
		Label:    label,
		Secret:   secret,
		Digits:   6,
		TimeStep: 30,
	}
}
