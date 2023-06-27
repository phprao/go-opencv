package util

import (
	"fmt"

	"gocv.io/x/gocv"
)

func ShowImage(title string, img gocv.Mat, shouldWaitKey bool) {
	w := gocv.NewWindow(title)
	w.ResizeWindow(500, 300)
	w.IMShow(img)
	if shouldWaitKey {
		w.WaitKey(0)
	}
}

func ReadAndShowImage(w *gocv.Window, filename string) gocv.Mat {
	img := gocv.IMRead(filename, gocv.IMReadColor)
	if img.Empty() {
		fmt.Printf("Error reading image from: %v\n", filename)
		return img
	}

	fmt.Println(img.Size())

	w.IMShow(img)
	return img
}

func ReadAndShowVideo(w *gocv.Window, v interface{}) {
	vc, err := gocv.OpenVideoCapture(v)
	if err != nil {
		fmt.Println(err)
		return
	}

	mat := gocv.NewMat()

	for {
		if vc.Read(&mat) {
			w.IMShow(mat)
			w.WaitKey(1000)
		} else {
			break
		}
	}

}
