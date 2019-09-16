package cui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"time"
)

type App struct {
	UI *gocui.Gui

	frames map[string]*AppFrame
}

func NewApp() (*App, error) {
	ui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, fmt.Errorf("new app cui failed: %s", err)
	}
	return &App{UI: ui, frames: make(map[string]*AppFrame)}, nil
}

func (a *App) AddToLayout(frame *AppFrame) {
	a.frames[frame.name] = frame
}

func (a *App) Run() error {
	a.UI.SetManagerFunc(func(gui *gocui.Gui) error {
		x, y := gui.Size()
		for name, frame := range a.frames {
			x0, y0, x1, y1 := frame.Position(x, y)
			view, err := gui.SetView(name, x0, y0, x1, y1)
			if err == gocui.ErrUnknownView {
				if frame.initCb != nil {
					err := frame.initCb(a, gui, view)
					if err != nil {
						return fmt.Errorf("init view[%s] failed: %s", name, err)
					}
				}
			} else if err != nil {
				return fmt.Errorf("generate view[%s] failed: %s", name, err)
			}

			if frame.updatedCb != nil {
				_ = frame.updatedCb(a, gui, view)
			}

		}

		return nil
	})

	err := a.bindKeys()
	if err != nil {
		return fmt.Errorf("bind keys failed; %s", err)
	}

	go func() {
		ticker := time.Tick(500 * time.Millisecond)
		for {
			select {
			case <-ticker:
				a.UI.Update(func(gui *gocui.Gui) error {
					return nil
				})
			}
		}
	}()

	return a.UI.MainLoop()
}

func (a *App) bindKeys() error {
	if err := a.UI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return fmt.Errorf("bind control-C failed; %s", err)
	}

	return nil
}

func (a *App) GetFrameByName(name string) (*AppFrame, bool) {
	f, ok := a.frames[name]
	return f, ok
}
