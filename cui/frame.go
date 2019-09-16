package cui

import (
	"bufio"
	"bytes"
	"github.com/jroimartin/gocui"
	"strings"
)

type ViewHandler func(app *App, g *gocui.Gui, v *gocui.View) error
type PositionSetter func(x, y int) (int, int, int, int)

type AppFrame struct {
	app *App

	name     string
	Position PositionSetter

	initCb    ViewHandler
	updatedCb ViewHandler

	view *gocui.View
}

func NewAppFrame(app *App, name string, position PositionSetter) (*AppFrame) {
	frame := &AppFrame{
		app: app, name: name, Position: position,
	}
	return frame
}

func (a *AppFrame) Init(cb ViewHandler) {
	a.initCb = cb
}

func (a *AppFrame) OnUpdated(cb ViewHandler) {
	a.updatedCb = cb
}

type LinesFrame struct {
	*AppFrame
}

func ShowLastBufferLines(buffer *bytes.Buffer) ViewHandler {
	return ShowLastBufferLinesWithSizedBuffer(buffer, -1)
}

func consumeBytesBuffer(buffer *bytes.Buffer, size int) {
	for i := 0; i < size; i++ {
		_, _ = buffer.ReadByte()
	}
}

func ShowLastBufferLinesWithSizedBuffer(buffer *bytes.Buffer, size int) ViewHandler {
	return func(app *App, g *gocui.Gui, v *gocui.View) error {
		buf := buffer.Bytes()

		if len(buf) > size && size > 0 {
			consumeBytesBuffer(buffer, len(buf)-size)
		}

		var (
			lines []string
		)

		lineScanner := bufio.NewScanner(bytes.NewBuffer(buf))
		lineScanner.Split(bufio.ScanLines)
		for lineScanner.Scan() {
			lines = append(lines, lineScanner.Text())
		}

		_, height := v.Size()

		if height > 1 && len(lines) > height {
			lines = lines[len(lines)-height:]
		}

		v.Clear()
		_, _ = v.Write([]byte(strings.Join(lines, "\n")))
		return nil
	}
}
