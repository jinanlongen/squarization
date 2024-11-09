package main

import (
	"image"
	utils "squarization/pkg/utils"
)

func main() {
	r1 := image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(20, 10)}
	r2 := image.Rectangle{Min: image.Pt(10, 10), Max: image.Pt(30, 20)}
	println("r1, r2:", r1.String(), r2.String())

	println("r1 eq r2:", r1.Eq(r2))
	println("r1.Intersect(r2):", r1.Intersect(r2).String())
	println("r1.Union(r2):", r1.Union(r2).String())
	println("r1.Add(image.Pt(20, 20)):", r1.Add(image.Pt(20, 20)).String())
	println("r1.Add(image.Pt(10, 10)):", r1.Add(image.Pt(10, 10)).Eq(r2))

	println("Min(1, 2, 3):", utils.Min(1, 2, 3))

	println("Center of[r1, r2]: ", r1.String(), utils.Center(r1).String(), r2.String(), utils.Center(r2).String())

	println("Max(1,2):", utils.Max(1, 2))

	println(utils.RectToSquare(r1).String())
}
