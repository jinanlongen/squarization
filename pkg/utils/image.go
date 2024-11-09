package utils

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func Center(rect image.Rectangle) image.Point {
	return image.Pt(rect.Min.X+rect.Dx()/2, rect.Min.Y+rect.Dy()/2)
}

func MaxCenteredSquareInRectangle(rect image.Rectangle, center image.Point) image.Rectangle {
	halfSquareSize := Min(
		center.X-rect.Min.X,
		rect.Max.X-center.X,
		center.Y-rect.Min.Y,
		rect.Max.Y-center.Y,
	)

	return image.Rect(
		center.X-halfSquareSize,
		center.Y-halfSquareSize,
		center.X+halfSquareSize,
		center.Y+halfSquareSize,
	)
}

func RectToSquare(rect image.Rectangle) image.Rectangle {
	sideLength := Max(rect.Dx(), rect.Dy())
	x0 := Center(rect).X - sideLength/2
	y0 := Center(rect).Y - sideLength/2
	return image.Rect(x0, y0, x0+sideLength, y0+sideLength)
}

func LargestContourRect(img gocv.Mat) image.Rectangle {
	// Convert to gray scale image
	grayImg := gocv.NewMat()
	defer grayImg.Close()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	threshImg := gocv.NewMat()
	defer threshImg.Close()
	gocv.Threshold(grayImg, &threshImg, 128, 255, gocv.ThresholdBinaryInv)

	contours := gocv.FindContours(threshImg, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	if contours.Size() == 0 {
		fmt.Println("Error: No contours found in the image.")
		return gocv.BoundingRect(gocv.PointVector{})
	}

	largestContour := contours.At(0)
	largestArea := gocv.ContourArea(largestContour)
	for i := 1; i < contours.Size(); i++ {
		contour := contours.At(i)
		area := gocv.ContourArea(contour)
		if area > largestArea {
			largestContour = contour
			largestArea = area
		}
	}
	return gocv.BoundingRect(largestContour)
}

func ContoursBoundingRect(img gocv.Mat) image.Rectangle {
	// Convert to gray scale image
	grayImg := gocv.NewMat()
	defer grayImg.Close()
	gocv.CvtColor(img, &grayImg, gocv.ColorBGRToGray)

	threshImg := gocv.NewMat()
	defer threshImg.Close()
	gocv.Threshold(grayImg, &threshImg, 128, 255, gocv.ThresholdBinaryInv)

	contours := gocv.FindContours(threshImg, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	if contours.Size() == 0 {
		fmt.Println("Error: No contours found in the image.")
		return gocv.BoundingRect(gocv.PointVector{})
	}

	contour := contours.At(0)
	rect := gocv.BoundingRect(contour)
	for i := 1; i < contours.Size(); i++ {
		contour = contours.At(i)
		rect = rect.Union(gocv.BoundingRect(contour))
	}

	return rect
}

func ShowMat(img gocv.Mat, title string) {
	window := gocv.NewWindow(title)
	defer window.Close()

	window.IMShow(img)
	gocv.WaitKey(0)
}

func CenterAndPadImage(canvasMat gocv.Mat, overlayMat gocv.Mat, paddingColor color.RGBA) image.Rectangle {
	canvasWidth := canvasMat.Cols()
	canvasHeight := canvasMat.Rows()

	overlayWidth := overlayMat.Cols()
	overlayHeight := overlayMat.Rows()

	hPadding := (canvasWidth - overlayWidth) / 2
	vPadding := (canvasHeight - overlayHeight) / 2

	gocv.CopyMakeBorder(overlayMat, &canvasMat, vPadding, vPadding, hPadding, hPadding, gocv.BorderConstant, paddingColor)

	overlayRect := image.Rect(hPadding, vPadding, canvasWidth-hPadding, canvasHeight-vPadding)
	return overlayRect
}

func BoundingRect(mat gocv.Mat) image.Rectangle {
	rows, cols := mat.Rows(), mat.Cols()
	return image.Rect(0, 0, cols, rows)
}

func MatToSquare(inMat gocv.Mat) (gocv.Mat, image.Rectangle) {
	boundingRectForInMat := BoundingRect(inMat)
	boundingSquareForInMat := RectToSquare(boundingRectForInMat)
	outMat := gocv.NewMatWithSize(boundingSquareForInMat.Dy(), boundingSquareForInMat.Dx(), gocv.MatTypeCV8UC3)

	overlayRect := CenterAndPadImage(outMat, inMat, GetBackgroundColor(inMat))

	return outMat, overlayRect
}

func Squarify(inMat gocv.Mat) gocv.Mat {
	contoursBoundingRect := LargestContourRect(inMat)
	// contoursBoundingRect := utils.ContoursBoundingRect(inMat)

	if contoursBoundingRect.Empty() {
		log.Println("Contours NOT found. ")

		canvasMat, overlayRect := MatToSquare(inMat)
		gocv.Rectangle(&canvasMat, overlayRect, color.RGBA{255, 0, 0, 0}, 2)

		return canvasMat
	} else {
		log.Println("Contours FOUND:", contoursBoundingRect.String())

		inMatRect := BoundingRect(inMat)
		contoursBoundingSquare := RectToSquare(contoursBoundingRect)
		unionRect := inMatRect.Union(contoursBoundingSquare)

		if unionRect.Eq(inMatRect) {
			log.Println("unionSquare, inMatRect: EQ", unionRect.String(), inMatRect.String())
			centerPt := Center(contoursBoundingRect)

			maxCenteredSquare := MaxCenteredSquareInRectangle(inMatRect, centerPt)
			gocv.Rectangle(&inMat, maxCenteredSquare, color.RGBA{255, 0, 0, 0}, 2)

			return inMat
		} else {
			log.Println("unionSquare, inMatRect: NOT EQ ", unionRect.String(), inMatRect.String())

			canvasMat, overlayRect := MatToSquare(inMat)
			gocv.Rectangle(&canvasMat, overlayRect, color.RGBA{255, 0, 0, 0}, 2)
			return canvasMat
		}
	}
}

func GetBackgroundColor(img gocv.Mat) color.RGBA {
	// Get dimensions
	rows := img.Rows()
	cols := img.Cols()

	// Sample the four corners
	topLeft := img.GetVecbAt(0, 0)
	topRight := img.GetVecbAt(0, cols-1)
	bottomLeft := img.GetVecbAt(rows-1, 0)
	bottomRight := img.GetVecbAt(rows-1, cols-1)

	// Sum the color values from the corners
	sumB := int(topLeft[0]) + int(topRight[0]) + int(bottomLeft[0]) + int(bottomRight[0])
	sumG := int(topLeft[1]) + int(topRight[1]) + int(bottomLeft[1]) + int(bottomRight[1])
	sumR := int(topLeft[2]) + int(topRight[2]) + int(bottomLeft[2]) + int(bottomRight[2])

	// Calculate the average color
	avgB := sumB / 4
	avgG := sumG / 4
	avgR := sumR / 4

	// Return the color as color.RGBA (with full opacity)
	return color.RGBA{R: uint8(avgR), G: uint8(avgG), B: uint8(avgB), A: 255}
}
