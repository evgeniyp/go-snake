package main

import (
	"fmt"
	"github.com/evgeniyp/go-snake/fonts"
	"image/color"
	"log"
	"math/rand"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	w         = 75
	h         = 50
	pixelSize = 8

	windowWidth  = pixelSize * w
	windowHeight = pixelSize * h

	fontPixelSize = 16
	fontSize      = 40
	fontDPI       = 72
	maxTPS        = 20
)

var (
	colorWhite     = color.Gray{0xff}
	colorRed       = color.RGBA{0xff, 0x0, 0x00, 0xff}
	colorGreen     = color.RGBA{0x00, 0xff, 0x00, 0xff}
	colorDarkGreen = color.RGBA{0x00, 0x80, 0x00, 0xff}
)

type Game struct {
	score      int
	dX, dY     int
	sizeToGrow int
	food       Coord
	snake      []Coord
	pixelFont  font.Face
	isRunning  bool
	gameImage  *ebiten.Image
}

type Coord struct {
	X int
	Y int
}

func (g *Game) Init() {
	g.score = 0
	g.dX = 1
	g.dY = 0
	g.food = Coord{rand.Intn(w), rand.Intn(h)}
	g.snake = []Coord{{1, 0}, {0, 0}}

	tt, _ := truetype.Parse(fonts.FontZxSpectrum7)
	g.pixelFont = truetype.NewFace(tt, &truetype.Options{
		Size:    fontSize,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})

	g.isRunning = true
	g.gameImage, _ = ebiten.NewImage(w, h, ebiten.FilterNearest)
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !g.isRunning {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Init()
		}
	} else {

		if g.dY == 0 && ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.dY = -1
			g.dX = 0
		} else if g.dY == 0 && ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.dY = +1
			g.dX = 0
		} else if g.dX == 0 && ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.dX = -1
			g.dY = 0
		} else if g.dX == 0 && ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.dX = +1
			g.dY = 0
		}

		newHead := Coord{g.snake[0].X + g.dX, g.snake[0].Y + g.dY}
		g.snake = append([]Coord{newHead}, g.snake...)

		if g.snake[0] == g.food {
			g.score++
			g.sizeToGrow++
			g.food.X = rand.Intn(w)
			g.food.Y = rand.Intn(h)
		}

		if g.sizeToGrow > 0 {
			g.sizeToGrow -= 1
		} else {
			g.snake = g.snake[:len(g.snake)-1]
		}

		if g.snake[0].X < 0 || g.snake[0].X > w-1 || g.snake[0].Y < 0 || g.snake[0].Y > h-1 {
			g.isRunning = false
		}
		for i := 1; i < len(g.snake); i++ {
			if g.snake[0] == g.snake[i] {
				g.isRunning = false
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gameImage.Fill(color.Black)
	g.gameImage.Set(g.food.X, g.food.Y, colorRed)
	for i, snakePart := range g.snake {
		var clr color.Color
		if i == 0 {
			clr = colorGreen
		} else {
			clr = colorDarkGreen
		}
		g.gameImage.Set(snakePart.X, snakePart.Y, clr)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(pixelSize), float64(pixelSize))

	screen.DrawImage(g.gameImage, op)
	var msg = fmt.Sprintf("Score: %v", g.score)
	text.Draw(screen, msg, g.pixelFont, pixelSize, (h-1)*pixelSize, color.White)

	if !g.isRunning {
		text.Draw(screen,
			"GAME  OVER",
			g.pixelFont,
			(windowWidth-len("GAME  OVER")*16)/2,
			windowHeight/2-fontPixelSize,
			colorRed)

		text.Draw(screen,
			"Press SPACE to Restart",
			g.pixelFont,
			(windowWidth-len("Press SPACE to Restart")*16)/2,
			windowHeight/2+fontPixelSize,
			colorWhite)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Snake")
	ebiten.SetWindowDecorated(true)
	ebiten.SetMaxTPS(maxTPS)
	g := &Game{}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
