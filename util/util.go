package util

import (
	"fmt"
	"image"

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

// 同时展示多个图片
func ShowMultipleImage(title string, imgs []gocv.Mat, imgCols int) {
	if imgs == nil {
		return
	}
	imgNum := len(imgs)
	imgOriSize := imgs[0].Size() // [行数 列数]
	imgDst := gocv.NewMatWithSize(imgOriSize[0]*((imgNum-1)/imgCols+1), imgOriSize[1]*imgCols, imgs[0].Type())
	imgChannel := 3 // 都转换成 BGR 通道

	m := gocv.NewMat()
	for i := 0; i < imgNum; i++ {
		// 像素点位置
		x0 := (i % imgCols) * imgOriSize[1]
		y0 := (i / imgCols) * imgOriSize[0]
		x1 := x0 + imgOriSize[1]
		y1 := y0 + imgOriSize[0]

		// Region 返回的 Mat 和原始的 Mat是引用关系，操作是相互影响的
		regin := imgDst.Region(image.Rect(x0, y0, x1, y1))
		if imgs[i].Channels() != imgChannel {
			gocv.CvtColor(imgs[i], &m, gocv.ColorGrayToBGR)
			m.CopyTo(&regin)
		} else {
			imgs[i].CopyTo(&regin)
		}
	}

	w := gocv.NewWindow(title)
	// imgDst 是一个整体，要求每一块的通道数一样，否则就不是一个合格的 Mat，无法展示。
	w.IMShow(imgDst)
	w.WaitKey(0)
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
	// gif, err := gocv.VideoCaptureFile("showimage/image15.gif")
	// gif, err := gocv.OpenVideoCapture("showimage/image15.gif")
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
