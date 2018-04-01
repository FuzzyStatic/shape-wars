/*
 * @Author: Allen Flickinger (allen.flickinger@gmail.com)
 * @Date: 2018-03-31 20:20:08
 * @Last Modified by: FuzzyStatic
 * @Last Modified time: 2018-03-31 20:20:49
 */

package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type player struct {
	lives            int64
	score            int64
	geometry         geometry
	missiles         []*geometry
	lastMissileFired time.Time
}

func (p *player) changePosition(win *pixelgl.Window) {
	if (win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW)) &&
		p.geometry.shape.max.y < win.Bounds().Max.Y {
		fmt.Println("up")
		p.geometry.shape.min.y = p.geometry.shape.min.y + p.geometry.velocity.speed
		p.geometry.shape.max.y = p.geometry.shape.max.y + p.geometry.velocity.speed
	}

	if (win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS)) &&
		p.geometry.shape.min.y > win.Bounds().Min.Y {
		fmt.Println("down")
		p.geometry.shape.min.y = p.geometry.shape.min.y - p.geometry.velocity.speed
		p.geometry.shape.max.y = p.geometry.shape.max.y - p.geometry.velocity.speed
	}

	if (win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA)) &&
		p.geometry.shape.min.x > win.Bounds().Min.X {
		fmt.Println("left")
		p.geometry.shape.min.x = p.geometry.shape.min.x - p.geometry.velocity.speed
		p.geometry.shape.max.x = p.geometry.shape.max.x - p.geometry.velocity.speed
	}

	if (win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD)) &&
		p.geometry.shape.max.x < win.Bounds().Max.X {
		fmt.Println("right")
		p.geometry.shape.min.x = p.geometry.shape.min.x + p.geometry.velocity.speed
		p.geometry.shape.max.x = p.geometry.shape.max.x + p.geometry.velocity.speed
	}

	updatePts(&p.geometry.shape)

	return
}

func (p *player) missileManager(win *pixelgl.Window) {
	if (win.Pressed(pixelgl.KeySpace) || win.Pressed(pixelgl.MouseButton1)) &&
		time.Now().Sub(p.lastMissileFired) > (time.Millisecond*200) {
		fmt.Println("space")
		var (
			missile geometry
			mVec    pixel.Vec
			cVec    pixel.Vec
		)

		mVec = win.MousePosition()
		cVec = getCentroid(p.geometry.shape.pts)

		missile = geometry{
			velocity: velocity{
				speed:     5.0,
				direction: math.Atan2(mVec.Y-cVec.Y, mVec.X-cVec.X),
			},
			shape: shape{
				imd:  imdraw.New(nil),
				size: 20,
				clr: map[int]pixel.RGBA{
					1: pixel.RGB(0, 0, 1),
					2: pixel.RGB(0, 0, 1),
					3: pixel.RGB(0, 0, 1),
					4: pixel.RGB(0, 0, 1),
				},
			},
		}
		missile.shape.min = vec{
			x: cVec.X - missile.shape.size/2,
			y: cVec.Y - missile.shape.size/2,
		}
		missile.shape.max = vec{
			x: cVec.X + missile.shape.size/2,
			y: cVec.Y + missile.shape.size/2,
		}
		updatePts(&missile.shape)
		missile.shape.imd.Polygon(0)

		p.lastMissileFired = time.Now()
		p.missiles = append(p.missiles, &missile)
		playWavAudio(WAVMISSILE)
	}

	for _, m := range p.missiles {
		m.shape.imd.Clear()
		changePosition(m)
		updatePosition(&m.shape)
		m.shape.imd.Draw(win)
	}

	return
}

func (p *player) missileCleanup(win *pixelgl.Window) {
	for i := 0; i < len(p.missiles); i++ {
		// Check if missile is out of bounds
		if p.missiles[i].shape.min.x > win.Bounds().Max.X ||
			p.missiles[i].shape.max.x < win.Bounds().Min.X ||
			p.missiles[i].shape.min.y > win.Bounds().Max.Y ||
			p.missiles[i].shape.max.y < win.Bounds().Min.Y {
			// Delete with preserved order
			copy(p.missiles[i:], p.missiles[i+1:])
			p.missiles[len(p.missiles)-1] = nil // or the zero value of T
			p.missiles = p.missiles[:len(p.missiles)-1]
		}
	}

	return
}

func (p *player) missileEnemyCollision(e *enemyList) {
	for i := 0; i < len(p.missiles); i++ {
		for j := 0; j < len(e.enemies); j++ {
			if rectanglesIntersect(p.missiles[i].shape.min, p.missiles[i].shape.max,
				e.enemies[j].shape.min, e.enemies[j].shape.max) {
				copy(p.missiles[i:], p.missiles[i+1:])
				p.missiles[len(p.missiles)-1] = nil // or the zero value of T
				p.missiles = p.missiles[:len(p.missiles)-1]

				copy(e.enemies[j:], e.enemies[j+1:])
				e.enemies[len(e.enemies)-1] = nil // or the zero value of T
				e.enemies = e.enemies[:len(e.enemies)-1]

				playWavAudio(WAVEXPLOSION)

				p.score++
				if math.Mod(float64(p.score), 5) == 0 {
					e.speedMod = e.speedMod + 0.05
				}

				if math.Mod(float64(p.score), 30) == 0 {
					e.spawnRate = e.spawnRate - e.spawnMod
				}

				if math.Mod(float64(p.score), 50) == 0 {
					p.lives++
				}

				p.missileEnemyCollision(e)
				break
			}
		}
	}

	return
}

func (p *player) playerEnemyCollision(e *enemyList) {
	for i := 0; i < len(e.enemies); i++ {
		if rectanglesIntersect(p.geometry.shape.min, p.geometry.shape.max,
			e.enemies[i].shape.min, e.enemies[i].shape.max) {
			copy(e.enemies[i:], e.enemies[i+1:])
			e.enemies[len(e.enemies)-1] = nil // or the zero value of T
			e.enemies = e.enemies[:len(e.enemies)-1]

			p.lives--
		}
	}

	return
}
