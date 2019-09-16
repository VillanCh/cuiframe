package main

import (
	"bytes"
	"cuiframe/cui"
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	app, err := cui.NewApp()
	if err != nil {
		logrus.Fatalln(err)
	}

	var (
		resultsBuffer  = bytes.NewBuffer([]byte{})
		logFrameBuffer = bytes.NewBuffer([]byte{})
	)

	// 添加 ProcessBar
	processBar := cui.NewAppFrame(app, "processbar", func(x, y int) (i int, i2 int, i3 int, i4 int) {
		return 0, 0, x - 1, 2
	})
	processBar.Init(func(app *cui.App, g *gocui.Gui, v *gocui.View) error {
		v.Title = "Progress Bar"
		return nil
	})
	app.AddToLayout(processBar)

	// 添加 Log Frame
	logFrame := cui.NewAppFrame(app, "log", func(x, y int) (i int, i2 int, i3 int, i4 int) {
		return 0, 3, x/3 - 1, y - 1
	})
	logFrame.Init(func(app *cui.App, g *gocui.Gui, v *gocui.View) error {
		v.Title = "Log View"
		return nil
	})
	logFrame.OnUpdated(cui.ShowLastBufferLinesWithSizedBuffer(logFrameBuffer, 500))
	app.AddToLayout(logFrame)

	results := cui.NewAppFrame(app, "results", func(x, y int) (i int, i2 int, i3 int, i4 int) {
		return x/3 + 1, 3, x - 1, y - 1
	})
	results.Init(func(app *cui.App, g *gocui.Gui, v *gocui.View) error {
		v.Title = "Results"
		return nil
	})
	results.OnUpdated(cui.ShowLastBufferLines(resultsBuffer))
	app.AddToLayout(results)

	// 每 0.5 秒写一次 Log
	go func() {
		ticker := time.Tick(100 * time.Millisecond)
		for {
			select {
			case <-ticker:
				resultsBuffer.WriteString(fmt.Sprintf("strinasdfasdfasdf %s full: %v \r\n", time.Now().String(), len(resultsBuffer.String())))
				logFrameBuffer.WriteString(fmt.Sprintf("logger is sized to 200 now: %s now: %v\n", time.Now().String(), len(logFrameBuffer.String())))
			}
		}
	}()

	if err := app.Run(); err != nil {
		logrus.Fatalf("run main loop failed: %s", err)
	}
}
