// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"log"
)

const (
	IDI_ICON = 11
)

func ToggleStyle(hwnd win.HWND, b bool, style int) {
	originalStyle := int(win.GetWindowLongPtr(hwnd, win.GWL_STYLE))
	if originalStyle != 0 {
		if b {
			originalStyle |= style
		} else {
			originalStyle ^= style
		}
		win.SetWindowLongPtr(hwnd, win.GWL_STYLE, uintptr(originalStyle))
	}
}

func EnableMaxButton(hwnd win.HWND, b bool) {
	ToggleStyle(hwnd, b, win.WS_MAXIMIZEBOX)
}

func EnableMinButton(hwnd win.HWND, b bool) {
	ToggleStyle(hwnd, b, win.WS_MINIMIZEBOX)
}

func EnableSizable(hwnd win.HWND, b bool) {
	ToggleStyle(hwnd, b, win.WS_THICKFRAME)
}

func Center(this *walk.MainWindow) {
	sWidth := win.GetSystemMetrics(win.SM_CXFULLSCREEN)
	sHeight := win.GetSystemMetrics(win.SM_CYFULLSCREEN)
	if sWidth != 0 && sHeight != 0 {
		size := this.Size()
		this.SetX(int(sWidth/2) - (size.Width / 2))
		this.SetY(int(sHeight/2) - (size.Height / 2))
	}
}

func main() {
	var mw *walk.MainWindow
	var edtRemoteAddr, edtRemotePort, edtLocalAddr, edtLocalPort *walk.LineEdit
	var edtLogView *walk.TextEdit
	var btnStart, btnStop *walk.PushButton
	var transfer *Transfer
	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Tcp Data Transfer  By:DingQi [Golang]",
		MinSize:  Size{450, 600},

		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 4, MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "RemoteAddr:",
					},
					LineEdit{
						AssignTo: &edtRemoteAddr,
						Text:     "10.34.80.1",
					},
					Label{
						Text: "RemotePort:",
					},
					LineEdit{
						AssignTo: &edtRemotePort,
						Text:     "8020",
					},
					Label{
						Text: "LocalAddr:",
					},
					LineEdit{
						AssignTo: &edtLocalAddr,
						Text:     "0.0.0.0",
					},
					Label{
						Text: "LocalPort:",
					},
					LineEdit{
						AssignTo: &edtLocalPort,
						Text:     "8018",
					},

					Label{
						ColumnSpan: 4,
						Text:       "LOG:",
					},

					TextEdit{
						AssignTo:   &edtLogView,
						ColumnSpan: 4,
						MinSize:    Size{100, 50},
					},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						AssignTo: &btnStart,
						Text:     "Start",
						OnClicked: func() {
							btnStop.SetEnabled(true)
							btnStart.SetEnabled(false)
							transfer = NewTransfer(edtRemoteAddr.Text(),
								edtRemotePort.Text(),
								edtLocalAddr.Text(),
								edtLocalPort.Text(),
							)
							transfer.Start()
						},
					},
					PushButton{
						AssignTo: &btnStop,
						Enabled:  false,
						Text:     "Stop",
						OnClicked: func() {
							transfer.Stop()
							btnStart.SetEnabled(true)
							btnStop.SetEnabled(false)
						},
					},
				},
			},
		},
	}).Create(); err != nil {
		log.Fatal(err)
	}

	if ico, err := walk.NewIconFromResourceId(IDI_ICON); err == nil {
		mw.SetIcon(ico)
	}

	lv, err := NewLogView(edtLogView)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(lv)

	EnableMaxButton(mw.Handle(), false)
	EnableSizable(mw.Handle(), false)
	Center(mw)
	mw.Show()
	mw.Run()
}
