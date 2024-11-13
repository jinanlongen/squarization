package main

import (
	"log"
	"os"
	"path/filepath"

	"squarization/pkg/utils"

	"gocv.io/x/gocv"
)

func main() {
	processDir("data", "outputs/")
}

func processDir(inputDir string, outputDir string) {
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking the path:", err)
			return err
		}

		if !info.IsDir() {
			inMat := gocv.IMRead(path, gocv.IMReadColor)
			if inMat.Empty() {
				log.Println("Error: Image not found or unable to load:", path)
			}
			defer inMat.Close()

			log.Println("Processing:", path)

			outMat := utils.Squarify(inMat)

			// utils.ShowMat(outMat, "Output")
			gocv.IMWrite(outputDir+info.Name()+".jpg", outMat)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}
}
