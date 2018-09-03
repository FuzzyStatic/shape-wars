/*
 * @Author: Allen Flickinger (allen.flickinger@gmail.com)
 * @Date: 2018-03-17 12:09:33
 * @Last Modified by: FuzzyStatic
 * @Last Modified time: 2018-03-31 20:28:41
 */

package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type hud struct {
	active bool
	text   *text.Text
}

func run() {
	var (
		win       *pixelgl.Window
		p         player
		e         enemyList
		ba        *text.Atlas
		debugHud  hud
		playerHud hud
		last      time.Time
		err       error
	)

	win, err = pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  "Shape Wars",
			Bounds: pixel.R(0, 0, 1024, 768),
			VSync:  true,
		})
	if err != nil {
		panic(err)
	}

	p = player{
		lives: 3,
		score: 0,
		geometry: geometry{
			velocity: velocity{
				speed: 5.0,
			},
			shape: shape{
				imd:  imdraw.New(nil),
				size: 50,
				clr: map[int]pixel.RGBA{
					1: pixel.RGB(0, 1, 0),
					2: pixel.RGB(0, 1, 0),
					3: pixel.RGB(0, 1, 0),
					4: pixel.RGB(0, 1, 0),
				},
			},
		},
	}
	p.geometry.shape.min = vec{
		x: win.Bounds().Max.X/2 - 50,
		y: win.Bounds().Max.Y/2 - 50,
	}
	p.geometry.shape.max = vec{
		x: win.Bounds().Max.X/2 + 50,
		y: win.Bounds().Max.Y/2 + 50,
	}
	updatePts(&p.geometry.shape)

	e = enemyList{
		lastCreatedEnemy: time.Now(),
		speed:            1,
		speedMod:         1,
		spawnRate:        time.Second * 1,
		spawnMod:         time.Millisecond * 100,
	}

	ba = text.NewAtlas(basicfont.Face7x13, text.ASCII)

	debugHud = hud{
		active: false,
		text:   text.New(pixel.V(10, 748), ba),
	}

	playerHud = hud{
		active: true,
		text:   text.New(pixel.V(1014, 748), ba),
	}

	win.Clear(colornames.Black)

	initAudo()

	last = time.Now()
	for !win.Closed() {
		var (
			m  runtime.MemStats
			dt float64
		)

		runtime.ReadMemStats(&m)

		dt = time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Black)

		p.geometry.shape.imd.Clear()
		p.changePosition(win)
		updatePosition(&p.geometry.shape)
		p.geometry.shape.imd.Draw(win)

		p.missileManager(win)
		e.enemyManager(win)

		p.missileEnemyCollision(&e)
		p.playerEnemyCollision(&e)

		p.missileCleanup(win)
		e.enemyCleanup(win)

		if win.JustPressed(pixelgl.KeyBackslash) {
			if debugHud.active {
				debugHud.active = false
			} else {
				debugHud.active = true
			}
		}

		debugHud.text.Clear()
		if debugHud.active {
			fmt.Fprintln(debugHud.text, "FPS: "+strconv.FormatFloat(1/dt, 'g', -1, 64))
			fmt.Fprintln(debugHud.text, "Alloc: "+strconv.FormatUint(m.Alloc/1024, 10))
			fmt.Fprintln(debugHud.text, "TotalAlloc: "+strconv.FormatUint(m.TotalAlloc/1024, 10))
			fmt.Fprintln(debugHud.text, "Sys: "+strconv.FormatUint(m.Sys/1024, 10))
			fmt.Fprintln(debugHud.text, "NumGC: "+strconv.FormatUint(uint64(m.NumGC), 10))
			debugHud.text.Draw(win, pixel.IM)
		}

		if win.JustPressed(pixelgl.KeySlash) {
			if playerHud.active {
				playerHud.active = false
			} else {
				playerHud.active = true
			}
		}

		playerHud.text.Clear()
		if playerHud.active {
			playerHud.text.Dot.X -= playerHud.text.BoundsOf("Lives: " + strconv.FormatInt(p.lives, 10)).W()
			fmt.Fprintln(playerHud.text, "Lives: "+strconv.FormatInt(p.lives, 10))
			playerHud.text.Dot.X -= playerHud.text.BoundsOf("Score: " + strconv.FormatInt(p.score, 10)).W()
			fmt.Fprintln(playerHud.text, "Score: "+strconv.FormatInt(p.score, 10))
			playerHud.text.Draw(win, pixel.IM)
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
