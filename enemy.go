/*
 * @Author: Allen Flickinger (allen.flickinger@gmail.com)
 * @Date: 2018-03-31 20:22:09
 * @Last Modified by: FuzzyStatic
 * @Last Modified time: 2018-03-31 20:23:19
 */

package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type enemyList struct {
	lastCreatedEnemy time.Time
	speed            float64
	speedMod         float64
	spawnRate        time.Duration
	spawnMod         time.Duration
	enemies          []*geometry
}

func (e *enemyList) enemyManager(win *pixelgl.Window) {
	if time.Now().Sub(e.lastCreatedEnemy) > (e.spawnRate) {
		e.enemies = append(e.enemies, e.createEnemy(win))
		e.lastCreatedEnemy = time.Now()
	}

	for _, e := range e.enemies {
		e.shape.imd.Clear()
		changePosition(e)
		updatePosition(&e.shape)
		e.shape.imd.Draw(win)
	}

	return
}

func (e *enemyList) createEnemy(win *pixelgl.Window) *geometry {
	var (
		enemy geometry
		r     int
	)

	enemy = geometry{
		velocity: velocity{
			speed: e.speed * e.speedMod,
		},
		shape: shape{
			imd:  imdraw.New(nil),
			size: 20,
			clr: map[int]pixel.RGBA{
				1: pixel.RGB(1, 0, 0),
				2: pixel.RGB(1, 0, 0),
				3: pixel.RGB(1, 0, 0),
				4: pixel.RGB(1, 0, 0),
			},
		},
	}

	rand.Seed(time.Now().UTC().UnixNano())
	r = rand.Intn(4)

	switch r {
	case 3: // If we get 3, let's random the x and start at the top
		enemy.shape.min.x = float64(rand.Intn(int(win.Bounds().Max.X)-
			int(win.Bounds().Min.X-enemy.shape.size)) + int(win.Bounds().Min.X-enemy.shape.size))
		enemy.shape.max.x = enemy.shape.min.x + enemy.shape.size
		enemy.shape.min.y = win.Bounds().Max.Y
		enemy.shape.max.y = win.Bounds().Max.Y + enemy.shape.size
		if enemy.shape.min.x < win.Bounds().Max.X/2 {
			enemy.velocity.direction = rand.Float64() * (-math.Pi / 2)
		} else {
			enemy.velocity.direction = rand.Float64()*((-math.Pi)-(-math.Pi/2)) + (-math.Pi / 2)
		}
		break
	case 2: // If we get 2, let's random the y and start at the left
		enemy.shape.min.x = win.Bounds().Min.X - enemy.shape.size
		enemy.shape.max.x = win.Bounds().Min.X
		enemy.shape.min.y = float64(rand.Intn(int(win.Bounds().Max.Y)-
			int(win.Bounds().Min.Y-enemy.shape.size)) + int(win.Bounds().Min.Y-enemy.shape.size))
		enemy.shape.max.y = enemy.shape.min.y + enemy.shape.size
		if enemy.shape.min.y < win.Bounds().Max.Y/2 {
			enemy.velocity.direction = rand.Float64() * (math.Pi / 2)
		} else {
			enemy.velocity.direction = rand.Float64() * (-math.Pi / 2)
		}
		break
	case 1: // If we get 1, let's random the x and start at the bottom
		enemy.shape.min.x = float64(rand.Intn(int(win.Bounds().Max.X)-
			int(win.Bounds().Min.X-enemy.shape.size)) + int(win.Bounds().Min.X-enemy.shape.size))
		enemy.shape.max.x = enemy.shape.min.x + enemy.shape.size
		enemy.shape.min.y = win.Bounds().Min.Y - enemy.shape.size
		enemy.shape.max.y = win.Bounds().Min.Y
		if enemy.shape.min.x < win.Bounds().Max.X/2 {
			enemy.velocity.direction = rand.Float64() * (math.Pi / 2)
		} else {
			enemy.velocity.direction = rand.Float64()*((math.Pi)-(math.Pi/2)) + (math.Pi / 2)
		}
		break
	case 0: // If we get 0, let's random the y and start at the right
		enemy.shape.min.x = win.Bounds().Max.X
		enemy.shape.max.x = win.Bounds().Max.X + enemy.shape.size
		enemy.shape.min.y = float64(rand.Intn(int(win.Bounds().Max.Y)-
			int(win.Bounds().Min.Y-enemy.shape.size)) + int(win.Bounds().Min.Y-enemy.shape.size))
		enemy.shape.max.y = enemy.shape.min.y + enemy.shape.size
		if enemy.shape.min.y < win.Bounds().Max.Y/2 {
			enemy.velocity.direction = rand.Float64()*((math.Pi)-(math.Pi/2)) + (math.Pi / 2)
		} else {
			enemy.velocity.direction = rand.Float64()*((-math.Pi)-(-math.Pi/2)) + (-math.Pi / 2)
		}
		break
	}
	updatePts(&enemy.shape)
	enemy.shape.imd.Polygon(0)

	return &enemy
}

func (e *enemyList) enemyCleanup(win *pixelgl.Window) {
	for i := 0; i < len(e.enemies); i++ {
		// Check if enemy is out of bounds. TODO: Make into function
		if e.enemies[i].shape.min.x > win.Bounds().Max.X+e.enemies[i].shape.size ||
			e.enemies[i].shape.max.x < win.Bounds().Min.X-e.enemies[i].shape.size ||
			e.enemies[i].shape.min.y > win.Bounds().Max.Y+e.enemies[i].shape.size ||
			e.enemies[i].shape.max.y < win.Bounds().Min.Y-e.enemies[i].shape.size {
			// Delete with preserved order
			copy(e.enemies[i:], e.enemies[i+1:])
			e.enemies[len(e.enemies)-1] = nil // or the zero value of T
			e.enemies = e.enemies[:len(e.enemies)-1]
		}
	}

	return
}
