package main

import (
	"fmt"
	"github.com/kagami/go-face"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

var facesDir = "uploads"
var rec *face.Recognizer

func StartRecognition()  {

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatalf("error opening web cam: %v", err)
	}
	defer webcam.Close()

	img := gocv.NewMat()
	defer img.Close()

	harrcascade := "/Users/thushara/go/src/gocv.io/x/gocv/data/haarcascade_frontalface_default.xml"
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load(harrcascade) {
		fmt.Println("Error reading cascade file: data/haarcascade_frontalface_default.xml")
		return
	}

	classifier.Load(harrcascade)
	defer classifier.Close()

	err = UpdateRecogniser()

	if err == nil {
		for {
			if ok := webcam.Read(&img); !ok || img.Empty() {
				log.Println("Unable to read from the webcam")
				continue
			}

			rects := classifier.DetectMultiScale(img)

			for _, r := range rects {
				imgFace := img.Region(r)
				buf, err := gocv.IMEncode(".jpg", imgFace)
				err = imgFace.Close()

				if err != nil {
					log.Println("Error Encode")
					continue
				}

				text := "Not Matched"

				faces, err := rec.Recognize(buf)

				if err == nil {
					for _, faced := range faces {
						faceId := rec.ClassifyThreshold(faced.Descriptor, float32(0.3))

						if faceId >= 0 {
							text = "Matched UserId - " + strconv.Itoa(faceId)
							_ = SaveAttendance(faceId, time.Now())
						}
					}
				}

				if err != nil {
					log.Println("Error Check Faces")
					log.Println(faces)
					continue
				}

				size := gocv.GetTextSize(text, gocv.FontHersheyComplexSmall, 3, 2)
				pt := image.Pt(r.Min.X + (r.Min.X/2) - size.X/2, r.Min.Y - 2)
				blue := color.RGBA{0,0,255,0}

				gocv.PutText(&img, text, pt, gocv.FontHersheyComplexSmall, 3, blue, 2)
				gocv.Rectangle(&img, r, blue, 2)

				fmt.Println(text)
			}
		}
	}
}

func UpdateRecogniser() error {

	newRec, err := face.NewRecognizer(facesDir)

	var faces []Face
	var samples []face.Descriptor
	var faceIds []int32

	result := GetAllFaces()
	faces = result.([]Face)

	for _, face := range faces {
		file := filepath.Join(facesDir, face.FilePath)
		user, err := newRec.RecognizeSingleFile(file)

		if err == nil && user != nil {
			samples = append(samples, user.Descriptor)
			faceIds = append(faceIds, int32(face.ProfileId))
		}
	}
	newRec.SetSamples(samples, faceIds)
	rec = newRec
	return err
}