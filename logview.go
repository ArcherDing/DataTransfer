package main

import (
	"errors"
	"syscall"
	"unsafe"
)

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type LogView struct {
	textEdit *walk.TextEdit
	logChan  chan string
}

const (
	TEM_APPENDTEXT = win.WM_USER + 6
)

func NewLogView(textEdit *walk.TextEdit) (*LogView, error) {
	lc := make(chan string, 1024)
	this := &LogView{textEdit: textEdit, logChan: lc}
	this.setReadOnly(true)
	this.textEdit.SendMessage(win.EM_SETLIMITTEXT, 4294967295, 0)
	return this, nil
}

func (this *LogView) setTextSelection(start, end int) {
	this.textEdit.SendMessage(win.EM_SETSEL, uintptr(start), uintptr(end))
}

func (this *LogView) textLength() int {
	return int(this.textEdit.SendMessage(0x000E, uintptr(0), uintptr(0)))
}

func (this *LogView) AppendText(value string) {
	textLength := this.textLength()
	this.setTextSelection(textLength, textLength)
	this.textEdit.SendMessage(win.EM_REPLACESEL, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))))
}

func (this *LogView) setReadOnly(readOnly bool) error {
	if 0 == this.textEdit.SendMessage(win.EM_SETREADONLY, uintptr(win.BoolToBOOL(readOnly)), 0) {
		return errors.New("fail to call EM_SETREADONLY")
	}

	return nil
}

func (this *LogView) GoAppendText(value string) {
	this.logChan <- value
	go this.MsgThread(TEM_APPENDTEXT)
}

func (this *LogView) Write(p []byte) (int, error) {
	this.GoAppendText(string(p) + "\r\n")
	return len(p), nil
}

func (this *LogView) MsgThread(msg uint32) {
	switch msg {
	case TEM_APPENDTEXT:
		select {
		case value := <-this.logChan:
			this.AppendText(value)
		default:
			return
		}
	}
}
