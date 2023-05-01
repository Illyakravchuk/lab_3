package ui

import (
	"image"
	"image/color"
	"log"
	"fmt"
	
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
	mousePos image.Point
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	pw.mousePos = image.Point{X:400, Y: 400}
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title: pw.Title,
		Width: 800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			// TODO: Реалізувати реакцію на натискання кнопки миші.
			if e.Button == mouse.ButtonLeft && e.Direction == mouse.DirPress {
				pw.mousePos = image.Point{
					X: int(e.X),
					Y: int(e.Y),
				}
				pw.w.Send(paint.Event{})
				fmt.Println("X pos, Y pos: ", e.X, e.Y)
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.White, draw.Src) // Фон.

	x, y := pw.mousePos.X, pw.mousePos.Y

	height, width := 400, 150
	yellow := color.RGBA{255, 255, 0, 255}

	pw.w.Fill(image.Rect(x-width/2, y+width, x+width/2, y), yellow, draw.Src)
	pw.w.Fill(image.Rect(x-height/2, y-height/2, x+height/2, y), yellow, draw.Src)

	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}