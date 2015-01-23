package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lazywei/go-opencv/opencv"
	"github.com/muesli/smartcrop"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
)

func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	width, _ := v["width"]
	height, _ := v["height"]
	imageUrl, _ := v["imageUrl"]

	x, _ := strconv.Atoi(width)
	y, _ := strconv.Atoi(height)

	// get the requested image
	res, err := http.Get("https://" + imageUrl)
	if err != nil || res.StatusCode != 200 {
		panic(err)
	}
	defer res.Body.Close()

	// load as an image.Image
	m, err := jpeg.Decode(res.Body)
	if err != nil {
		panic(err)
	}

	srcWidth := m.Bounds().Size().X
	srcHeight := m.Bounds().Size().Y
	newWidth := srcWidth
	newHeight := srcHeight
	srcRatio := float32(srcWidth) / float32(srcHeight)
	targetRatio := float32(x) / float32(y)
	// smart crop
	if srcRatio > targetRatio {
		// height stays the same
		newHeight = int(float32(srcWidth) / targetRatio)
	} else {
		// width stays the same
		newWidth = int(float32(srcHeight) * targetRatio)
	}
	topCrop, _ := smartcrop.SmartCrop(&m, newWidth, newHeight)
	// convert to opencv image
	srcImage := opencv.FromImage(m)
	if srcImage == nil {
		fmt.Printf("Couldn't create opencv Image")
	}
	defer srcImage.Release()

	// resize using x/y dimensions
	// constrain proportions by cropping
	croppedImage := opencv.Crop(srcImage, topCrop.X, topCrop.Y, topCrop.Width, topCrop.Height)
	resizedImage := opencv.Resize(croppedImage, x, y, opencv.CV_INTER_LINEAR)
	jpeg.Encode(w, resizedImage.ToImage(), &jpeg.Options{Quality: 90})
}

func main() {
	fmt.Printf("Setting up routes ... \n")
	r := mux.NewRouter()
	r.HandleFunc("/{width:[0-9]+}/{height:[0-9]+}/{imageUrl:.+}", ResizeHandler)
	http.Handle("/", r)

	fmt.Printf("Starting http Server ...\n")
	err := http.ListenAndServe("0.0.0.0:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
