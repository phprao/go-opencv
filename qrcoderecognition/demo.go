package qrcoderecognition

// 识别二维码

// opencv 自带的qrcode检测识别，其对图片的要求较高，因此识别精准度不高。

// https://blog.51cto.com/jsxyhelu2017/5972864
// https://blog.csdn.net/Yong_Qi2015/article/details/107194439

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/phprao/go-opencv/util"
	"gocv.io/x/gocv"
)

func Run() {
	src := gocv.IMRead("qrcoderecognition/1.jpg", gocv.IMReadColor)

	// 灰度
	gray := gocv.NewMat()
	gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)

	// 二值化
	threshold_output := gocv.NewMat()
	gocv.Threshold(gray, &threshold_output, 112, 255, gocv.ThresholdBinary)

	// 捕捉轮廓
	//
	// topology-拓扑结构
	// hierarchy内每个元素的4个int型变量是hierarchy[i][0] ~ hierarchy[i][3]，分别表示当前轮廓 i 的后一个轮廓、前一个轮廓、第一个子轮廓和父轮廓的编号索引。
	// 编号从0开始，如果当前轮廓没有对应的这四个关系轮廓，则相应的hierarchy[i][*]被置为-1。
	// 矩阵类型为CV32SC4，也就是每一个元素有四个分量，即后前子父
	// 在gocv中，hierarchy 是 1 行 n 列，具体使用参考 gocv 的 imgproc_test.go
	hierarchy := gocv.NewMat()
	pointsVector := gocv.FindContoursWithParams(threshold_output, &hierarchy, gocv.RetrievalTree, gocv.ChainApproxNone)

	// fmt.Println(hierarchy.Size(), hierarchy.Type()) // [1 967] CV32SC4
	// fmt.Println(pointsVector.Size())                // 967
	if pointsVector.IsNil() {
		log.Fatal("FindContoursWithParams is nil")
	}
	if pointsVector.Size() != hierarchy.Cols() {
		log.Fatal("FindContoursWithParams error")
	}

	drawSrc := gocv.Zeros(src.Rows(), src.Cols(), src.Type())
	src.CopyTo(&drawSrc)

	// 寻找子轮廓：三个角的定位块，其构造为：黑，白，黑，也就是三个回字形的轮廓
	levelContour := make([]int, 0)
	for i := 0; i < pointsVector.Size(); i++ {
		tempindex1 := hierarchy.GetVeciAt(0, i)[2]
		if tempindex1 != -1 {
			tempindex2 := hierarchy.GetVeciAt(0, int(tempindex1))[2]
			if tempindex2 != -1 {
				// 有些轮廓太小，可能就是个点，因此过滤掉面积很小的轮廓
				firstArea := gocv.ContourArea(pointsVector.At(i)) / gocv.ContourArea(pointsVector.At(int(tempindex1)))
				secondArea := gocv.ContourArea(pointsVector.At(int(tempindex1))) / gocv.ContourArea(pointsVector.At(int(tempindex2)))
				if (firstArea > 1 && firstArea < 10) && (secondArea > 1 && secondArea < 10) {
					levelContour = append(levelContour, i, int(tempindex1), int(tempindex2))
				}
			}
		}

		// 不停地按 1，可以看到整个画画的过程，用于调试
		// gocv.DrawContours(&drawSrc, pointsVector, i, color.RGBA{255, 0, 0, 255}, 1)
		// util.ShowImage("qrcode", drawSrc, true)
	}

	// fmt.Println(levelContour) //[230 231 232 298 299 300 304 305 306]
	if len(levelContour) != 9 {
		// 异常情况，需要进一步修正 TODO
		log.Fatal("levelContour > 9")
	}

	// 遍历 levelContour 拿出第三级轮廓，即 2 的值
	// 找到其重心，即X和Y各自的平均值
	centerPoint := make([]image.Point, 0)
	for i := 0; i < len(levelContour); i += 3 {
		gocv.DrawContours(&drawSrc, pointsVector, levelContour[i], color.RGBA{255, 0, 0, 255}, -1)

		points := pointsVector.At(levelContour[i]).ToPoints()
		pointsLen := len(points)
		sumX := 0
		sumY := 0
		for _, p := range points {
			sumX += p.X
			sumY += p.Y
		}
		centerPoint = append(centerPoint, image.Point{sumX / pointsLen, sumY / pointsLen})
	}
	// fmt.Println(centerPoint)

	// gocv.Line(&drawSrc, centerPoint[0], centerPoint[1], color.RGBA{0, 255, 0, 255}, 1)
	// gocv.Line(&drawSrc, centerPoint[1], centerPoint[2], color.RGBA{0, 255, 0, 255}, 1)
	// gocv.Line(&drawSrc, centerPoint[2], centerPoint[0], color.RGBA{0, 255, 0, 255}, 1)

	// 找到距离最大的边
	len01 := math.Sqrt(math.Pow(float64(centerPoint[0].X-centerPoint[1].X), 2) + math.Pow(float64(centerPoint[0].Y-centerPoint[1].Y), 2))
	len02 := math.Sqrt(math.Pow(float64(centerPoint[0].X-centerPoint[2].X), 2) + math.Pow(float64(centerPoint[0].Y-centerPoint[2].Y), 2))
	len12 := math.Sqrt(math.Pow(float64(centerPoint[1].X-centerPoint[2].X), 2) + math.Pow(float64(centerPoint[1].Y-centerPoint[2].Y), 2))

	/*
		0  2
		1  3
	*/
	centerPointNew := make([]image.Point, 4)
	if len01 > len02 && len01 > len12 {
		centerPointNew[0] = centerPoint[2]
		if centerPoint[0].Y > centerPoint[1].Y {
			centerPointNew[1] = centerPoint[0]
			centerPointNew[2] = centerPoint[1]
		} else {
			centerPointNew[1] = centerPoint[1]
			centerPointNew[2] = centerPoint[0]
		}
	}
	if len02 > len01 && len02 > len12 {
		centerPointNew[0] = centerPoint[1]
		if centerPoint[0].Y > centerPoint[2].Y {
			centerPointNew[1] = centerPoint[0]
			centerPointNew[2] = centerPoint[2]
		} else {
			centerPointNew[1] = centerPoint[2]
			centerPointNew[2] = centerPoint[0]
		}
	}
	if len12 > len01 && len12 > len02 {
		centerPointNew[0] = centerPoint[0]
		if centerPoint[1].Y > centerPoint[2].Y {
			centerPointNew[1] = centerPoint[1]
			centerPointNew[2] = centerPoint[2]
		} else {
			centerPointNew[1] = centerPoint[2]
			centerPointNew[2] = centerPoint[1]
		}
	}
	// fmt.Println(centerPointNew)
	centerPointNew[3] = image.Point{
		centerPointNew[2].X - centerPointNew[0].X + centerPointNew[1].X,
		centerPointNew[2].Y - centerPointNew[0].Y + centerPointNew[1].Y,
	}

	// gocv.Line(&drawSrc, centerPointNew[0], centerPointNew[1], color.RGBA{0, 255, 0, 255}, 1)
	// gocv.Line(&drawSrc, centerPointNew[0], centerPointNew[2], color.RGBA{0, 255, 0, 255}, 1)
	// gocv.Line(&drawSrc, centerPointNew[3], centerPointNew[1], color.RGBA{0, 255, 0, 255}, 1)
	// gocv.Line(&drawSrc, centerPointNew[3], centerPointNew[2], color.RGBA{0, 255, 0, 255}, 1)

	// 手机拍照如果没有平行于二维码所在平面，排出的二维码就发生了透视投影，需要校正，且只需要二维码这一个区域
	temp := 50
	pointVectorBefore := gocv.NewPointVectorFromPoints(centerPointNew)
	pointVectorAfter := gocv.NewPointVectorFromPoints([]image.Point{
		{0 + temp, 0 + temp}, {0 + temp, 100 + temp},
		{100 + temp, 0 + temp}, {100 + temp, 100 + temp}})
	mt := gocv.GetPerspectiveTransform(pointVectorBefore, pointVectorAfter)
	projectionDst := gocv.NewMat()
	gocv.WarpPerspectiveWithParams(src, &projectionDst, mt, image.Point{200, 200}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
	fmt.Println(projectionDst.Size())

	// 有时候二维码在打印的时候被拉伸了也要能识别出来

	// 如果得到的个数超过3个，需要将多余的删掉，或者可能有多个二维码，需要计算定位块之间的角度是否接近90度

	// drawPointsVector := gocv.NewMatWithSize(gray.Rows(), gray.Cols(), gocv.MatTypeCV8UC3)
	// gocv.DrawContours(&drawPointsVector, pointsVector, -1, color.RGBA{255, 0, 0, 0}, 1)
	util.ShowImage("qrcode", projectionDst, true)
}
