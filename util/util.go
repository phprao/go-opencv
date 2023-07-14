package util

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/freetype"
	"gocv.io/x/gocv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func ShowImage(title string, img gocv.Mat, shouldWaitKey bool) {
	w := gocv.NewWindow(title)
	w.ResizeWindow(img.Cols(), img.Rows())
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
	imgChannel := imgs[0].Channels() // 都转换成 BGR 通道

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

func ReadAndShowVideo(filename string) {
	w := gocv.NewWindow(filename)
	vc, err := gocv.VideoCaptureFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	mat := gocv.NewMat()

	for {
		if vc.Read(&mat) {
			w.IMShow(mat)
			w.WaitKey(10)
		} else {
			break
		}
	}
	w.WaitKey(0)
}

func ReadAndShowGIF(filename string) {
	w := gocv.NewWindow(filename)

	f, _ := os.Open(filename)
	defer f.Close()

	gi, _ := gif.DecodeAll(f)

	for k, v := range gi.Image {
		img, err := gocv.ImageToMatRGB(v)
		if err != nil {
			log.Fatal(err)
		}

		w.IMShow(img)
		w.WaitKey(gi.Delay[k] * 10) // delay 单位是百分之一秒，waitkey参数为毫秒
	}

	w.WaitKey(0)
}

func ReadAndShowImageFromUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	m, err := gocv.IMDecode(by, gocv.IMReadColor)
	if err != nil {
		return
	}
	w := gocv.NewWindow("url image")
	w.IMShow(m)
	w.WaitKey(0)
}

func WriteTextOnMat(mat *gocv.Mat, text string, textPos image.Point, fontFile string, textSize float64, textColor color.RGBA) error {
	img, err := mat.ToImage()
	if err != nil {
		return err
	}
	img, err = WriteTextOnImage(img, text, textPos, fontFile, textSize, textColor)
	if err != nil {
		return err
	}
	mat2, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		return err
	}
	mat2.CopyTo(mat)
	mat2.Close()
	return nil
}

// Image2RGBA Image2RGBA
func Image2RGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}
	baseSrcBounds := img.Bounds().Max
	w, h := baseSrcBounds.X, baseSrcBounds.Y
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	//copy图片
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Over)

	return dst
}

func WriteTextOnImage(img image.Image, text string, textPos image.Point, fontFile string, textSize float64, textColor color.RGBA) (image.Image, error) {
	// Read the font data.
	dpi := 72.0
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		fmt.Println(err)
		return img, err
	}
	ft, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		return img, err
	}

	// Initialize the context.
	fontColor := image.NewUniform(textColor)
	rgbaImg := Image2RGBA(img)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(ft)
	c.SetFontSize(textSize)
	c.SetClip(img.Bounds())
	c.SetDst(rgbaImg)
	c.SetSrc(fontColor)

	// Draw the text.
	pt := freetype.Pt(textPos.X, textPos.Y+int(c.PointToFixed(textSize)>>6))

	_, err = c.DrawString(text, pt)
	if err != nil {
		fmt.Println(err)
		return img, err
	}
	return rgbaImg, nil
}

// 展示直方图，单列
func SaveHistSingle[V comparable](data []V, title string, name string, filename string) {
	barData := make([]opts.BarData, len(data))
	keys := make([]string, len(data))
	for k, v := range data {
		barData[k] = opts.BarData{Value: v}
		keys[k] = strconv.Itoa(k)
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))

	bar.SetXAxis(keys).AddSeries(name, barData)
	f, _ := os.Create(filename)
	bar.Render(f)
}

// 折线图，单列
func SaveLineSingle[V comparable](data []V, title string, name string, filename string, smooth bool) {
	lineData := make([]opts.LineData, len(data))
	keys := make([]string, len(data))
	for k, v := range data {
		lineData[k] = opts.LineData{Value: v}
		keys[k] = strconv.Itoa(k)
	}

	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))

	line.SetXAxis(keys).AddSeries(name, lineData).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: smooth}))
	f, _ := os.Create(filename)
	line.Render(f)
}

func GetRandomBytes(rows int, cols int) []byte {
	data := make([]byte, rows*cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			data[i*cols+j] = uint8(rand.Intn(255))
		}
	}

	return data
}
