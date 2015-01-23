package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"os"
	//"github.com/lazywei/go-opencv/opencv"
	"github.com/disintegration/imaging"
	"github.com/muesli/smartcrop"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
)

func getImg(url string) image.Image {
	// get the requested image
	res, err := http.Get("https://" + url)
	if err != nil || res.StatusCode != 200 {
		panic(err)
	}
	defer res.Body.Close()

	// load as an image.Image
	m, err := jpeg.Decode(res.Body)
	if err != nil {
		panic(err)
	}
	return m
}

func getFillCrop(src image.Image, targetRatio float32) (smartcrop.Crop, error) {
	srcWidth := src.Bounds().Size().X
	srcHeight := src.Bounds().Size().Y
	newWidth := srcWidth
	newHeight := srcHeight
	srcRatio := float32(srcWidth) / float32(srcHeight)
	if srcRatio > targetRatio {
		// width stays the same
		newHeight = int(float32(srcWidth) / targetRatio)
	}
	if srcRatio < targetRatio {
		// height stays the same
		newWidth = int(float32(srcHeight) * targetRatio)
	}

	return smartcrop.SmartCrop(&src, newWidth, newHeight)
}

func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	// parse url vars
	v := mux.Vars(r)
	width, _ := v["width"]
	height, _ := v["height"]
	x, _ := strconv.Atoi(width)
	y, _ := strconv.Atoi(height)

	imageUrl, _ := v["imageUrl"]
	m := getImg(imageUrl)

	cropBox, _ := getFillCrop(m, float32(x)/float32(y))
	cropRect := image.Rect(cropBox.X, cropBox.Y, cropBox.X+cropBox.Width, cropBox.Y+cropBox.Height)
	croppedImg := imaging.Crop(m, cropRect)

	// convert to opencv image
	//srcImage := opencv.FromImage(m)
	//if srcImage == nil {
	//  fmt.Printf("Couldn't create opencv Image")
	//}
	//defer srcImage.Release()
	//croppedImage := opencv.Crop(srcImage, cropBox.X, cropBox.Y, cropBox.Width, cropBox.Height)
	//resizedImage := opencv.Resize(croppedImage, x, y, opencv.CV_INTER_LINEAR)
	resizedImage := imaging.Resize(croppedImg, x, y, imaging.CatmullRom)
	jpeg.Encode(w, resizedImage, &jpeg.Options{Quality: 90})
}

func main() {
	fmt.Printf("Setting up routes ... \n")
	r := mux.NewRouter()
	r.HandleFunc("/{width:[0-9]+}/{height:[0-9]+}/{imageUrl:.+}", ResizeHandler)
	http.Handle("/", r)

	fmt.Printf("Starting http Server ...\n")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
