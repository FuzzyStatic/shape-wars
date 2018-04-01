/*
 * @Author: Allen Flickinger (allen.flickinger@gmail.com)
 * @Date: 2018-03-31 20:21:24
 * @Last Modified by: FuzzyStatic
 * @Last Modified time: 2018-03-31 20:23:34
 */

package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	geo "github.com/paulmach/go.geo"
)

type velocity struct {
	speed     float64
	direction float64
}

type vec struct {
	x float64
	y float64
}

type shape struct {
	imd  *imdraw.IMDraw
	size float64
	min  vec
	max  vec
	clr  map[int]pixel.RGBA
	pts  map[int]pixel.Vec
}

type geometry struct {
	velocity velocity
	shape    shape
}

func updatePts(shape *shape) {
	shape.pts = map[int]pixel.Vec{
		1: pixel.V(shape.min.x, shape.min.y),
		2: pixel.V(shape.max.x, shape.min.y),
		3: pixel.V(shape.max.x, shape.max.y),
		4: pixel.V(shape.min.x, shape.max.y),
	}

	return
}

func updatePosition(shape *shape) {
	for index, vector := range shape.pts {
		shape.imd.Color = shape.clr[index]
		shape.imd.Push(vector)
	}
	shape.imd.Polygon(0)

	return
}

func changePosition(g *geometry) {
	g.shape.min.x = g.shape.min.x + (math.Cos(g.velocity.direction) * g.velocity.speed)
	g.shape.min.y = g.shape.min.y + (math.Sin(g.velocity.direction) * g.velocity.speed)
	g.shape.max.x = g.shape.max.x + (math.Cos(g.velocity.direction) * g.velocity.speed)
	g.shape.max.y = g.shape.max.y + (math.Sin(g.velocity.direction) * g.velocity.speed)
	updatePts(&g.shape)

	return
}

func getCentroid(pts map[int]pixel.Vec) pixel.Vec {
	var (
		ps geo.PointSet
	)

	for _, vector := range pts {
		ps.Push(geo.NewPoint(vector.X, vector.Y))
	}

	return pixel.Vec{
		X: ps.Centroid().X(),
		Y: ps.Centroid().Y(),
	}
}

func rectanglesIntersect(r1Min, r1Max, r2Min, r2Max vec) bool {
	/* Find if two rectangles overlap
	*	                   r2Max
	*			  |--------|
	*             |	 r1Max |
	*		|-----|--|	   |
	*  		|     |__|_____|
	*		| r2Min  |
	*   	|________|
	*   r1Min
	 */
	if r1Min.x > r2Max.x || r2Min.x > r1Max.x {
		return false
	}

	if r1Max.y < r2Min.y || r2Max.y < r1Min.y {
		return false
	}

	return true
}
