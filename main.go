package main

import (
	"cuiframe/cui"
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"
	"github.com/sirupsen/logrus"
	"time"
)

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	app, err := cui.NewApp()
	if err != nil {
		logrus.Fatalln(err)
	}
	defer func() {
		app.UI.Close()
		termbox.SetCursor(1, 1)
		_ = termbox.Flush()
	}()

	app.UI.Mouse = true
	app.UI.Cursor = true
	app.UI.ASCII = true

	var (
		resultsBuffer  = logrus.New()
		logFrameBuffer = logrus.New()
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
		v.Wrap = true
		v.Autoscroll = true
		v.Editable = true
		logFrameBuffer.SetOutput(v)
		if _, err := g.SetCurrentView("log"); err != nil {
			return err
		}

		return nil
	})
	app.AddToLayout(logFrame)

	results := cui.NewAppFrame(app, "results", func(x, y int) (i int, i2 int, i3 int, i4 int) {
		return x/3 + 1, 3, x - 1, y - 1
	})
	results.Init(func(app *cui.App, g *gocui.Gui, v *gocui.View) error {
		v.Title = "Results"
		v.Wrap = true
		v.Editable = true
		v.Overwrite = true
		resultsBuffer.SetOutput(v)
		if _, err := g.SetCurrentView("results"); err != nil {
			return err
		}

		return nil
	})
	if err := app.UI.SetKeybinding("results", gocui.KeyArrowUp, gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) error {
			if err := scrollView(view, -1); err != nil {
				return fmt.Errorf("set results scroller up failed: %s", err)
			}
			return nil
		}); err != nil {
		return
	}
	if err := app.UI.SetKeybinding("results", gocui.KeyArrowDown, gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) error {
			if err := scrollView(view, 1); err != nil {
				return fmt.Errorf("set results scroller down failed: %s", err)
			}
			return nil
		}); err != nil {
		return
	}
	app.UI.SetKeybinding("results", gocui.MouseWheelDown, gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) error {
			return scrollView(view, 1)
		})
	app.UI.SetKeybinding("results", gocui.MouseWheelUp, gocui.ModNone,
		func(gui *gocui.Gui, view *gocui.View) error {
			return scrollView(view, -1)
		})
	app.AddToLayout(results)

	// 每 0.5 秒写一次 Log
	go func() {
		ticker := time.Tick(100 * time.Millisecond)
		for {
			select {
			case <-ticker:
				resultsBuffer.Info(fmt.Sprintf("strinasdfasdfasdf %s full: %v \r\n", time.Now().String()))
				logFrameBuffer.Info(fmt.Sprintf("logger is sized to 200 now: %s now: %v\n", time.Now().String()))
			}
		}
	}()

	if err := app.Run(); err != nil {
		logrus.Fatalf("run main loop failed: %s", err)
	}
}
