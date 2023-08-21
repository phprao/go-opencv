package qrcoderecognition

/*

从图像中识别二维码

https://blog.51cto.com/jsxyhelu2017/5972864
https://blog.csdn.net/Yong_Qi2015/article/details/107194439

相机不一定平行于二维码平面，此处需要做投影变换

即便相机平行于二维码平面，也有可能发生角度旋转。

首先找到三个定位标志，判断他们之间的关系
夹角是90度，视为相机平行于二维码屏幕；否则不平行。
1、平行
	是否发生旋转，即判断上边是否与X轴平行。
	宽高是否相等，
2、不平行
	存在投影，会导致二维码各个边并不是平行四边形（比如 1-1.jpg），计算难度加大
*/

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/phprao/go-opencv/util"
	"gocv.io/x/gocv"
)

func Run() {
	src := gocv.IMRead("qrcoderecognition/5.jpg", gocv.IMReadColor)

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
	// 矩阵类型为CV32SC4，也就是每一个元素有四个分量，即后、前、子、父
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
		if gocv.ContourArea(pointsVector.At(i)) < 1 {
			continue
		}
		son := hierarchy.GetVeciAt(0, i)[2]
		if son != -1 {
			sonson := hierarchy.GetVeciAt(0, int(son))[2]
			if sonson != -1 {
				fatherArea := gocv.ContourArea(pointsVector.At(i))
				sonArea := gocv.ContourArea(pointsVector.At(int(son)))
				sonsonArea := gocv.ContourArea(pointsVector.At(int(sonson)))
				fmt.Printf("fatherArea:%2f, sonArea:%2f, sonsonArea:%2f\n", fatherArea, sonArea, sonsonArea)
				// 注意，满足三层轮廓嵌套关系的有定位块(Position Detection Pattern)和对齐块(Alignment Pattern)，二维码内容多了之后就会出现对齐块
				// 区别在于他们的边长比，定位块为 7:5:3，对齐块为 5:3:1
				if sonArea/sonsonArea > 7 {
					continue
				}
				// 有些轮廓太小，可能就是个点，因此过滤掉面积很小的轮廓
				firstArea := fatherArea / sonArea
				secondArea := sonArea / sonsonArea
				if (firstArea > 1 && firstArea < 10) && (secondArea > 1 && secondArea < 10) {
					levelContour = append(levelContour, i, int(son), int(sonson))
				}
			}
		}

		// 不停地按 1，可以看到整个画画的过程，用于调试
		// gocv.DrawContours(&drawSrc, pointsVector, i, color.RGBA{255, 0, 0, 255}, 1)
		// util.ShowImage("qrcode", drawSrc, true)
	}

	// fmt.Println(levelContour) //[230 231 232 298 299 300 304 305 306]
	if len(levelContour) != 9 {
		fmt.Println(len(levelContour))
		// 异常情况，需要进一步修正 TODO
		for _, v := range levelContour {
			gocv.DrawContours(&drawSrc, pointsVector, int(v), color.RGBA{255, 0, 0, 255}, 1)
			util.ShowImage("qrcode", drawSrc, true)
		}

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

	len01 = math.Sqrt(math.Pow(float64(centerPointNew[0].X-centerPointNew[1].X), 2) + math.Pow(float64(centerPointNew[0].Y-centerPointNew[1].Y), 2))
	len02 = math.Sqrt(math.Pow(float64(centerPointNew[0].X-centerPointNew[2].X), 2) + math.Pow(float64(centerPointNew[0].Y-centerPointNew[2].Y), 2))
	len12 = math.Sqrt(math.Pow(float64(centerPointNew[1].X-centerPointNew[2].X), 2) + math.Pow(float64(centerPointNew[1].Y-centerPointNew[2].Y), 2))

	// 夹角是否为90正负2度
	angle := 180 - (math.Asin(len01/len12)+math.Asin(len02/len12))*180/math.Pi
	fmt.Println("angle", angle) // 89度

	// fmt.Println(centerPointNew)
	centerPointNew[3] = image.Point{
		centerPointNew[2].X - centerPointNew[0].X + centerPointNew[1].X,
		centerPointNew[2].Y - centerPointNew[0].Y + centerPointNew[1].Y,
	}

	gocv.Line(&drawSrc, centerPointNew[0], centerPointNew[1], color.RGBA{0, 255, 0, 255}, 1)
	gocv.Line(&drawSrc, centerPointNew[0], centerPointNew[2], color.RGBA{0, 255, 0, 255}, 1)
	gocv.Line(&drawSrc, centerPointNew[3], centerPointNew[1], color.RGBA{0, 255, 0, 255}, 1)
	gocv.Line(&drawSrc, centerPointNew[3], centerPointNew[2], color.RGBA{0, 255, 0, 255}, 1)
	// util.ShowImage("qrcode0", drawSrc, true)

	// 手机拍照如果没有平行于二维码所在平面，排出的二维码就发生了透视投影，需要校正，且只需要二维码这一个区域
	projectionDst := perspective(src, centerPointNew)

	// 有时候二维码在打印的时候被拉伸了也要能识别出来

	// 如果得到的个数超过3个，需要将多余的删掉，或者可能有多个二维码，需要计算定位块之间的角度是否接近90度

	// drawPointsVector := gocv.NewMatWithSize(gray.Rows(), gray.Cols(), gocv.MatTypeCV8UC3)
	// gocv.DrawContours(&drawPointsVector, pointsVector, -1, color.RGBA{255, 0, 0, 0}, 1)

	content, _ := DetectAndDecodeQrcode(projectionDst)
	fmt.Println(content)

	util.ShowImage("qrcode", projectionDst, true)
}

func perspective(src gocv.Mat, pts []image.Point) gocv.Mat {
	temp := 50
	pointVectorBefore := gocv.NewPointVectorFromPoints(pts)
	pointVectorAfter := gocv.NewPointVectorFromPoints([]image.Point{
		{0 + temp, 0 + temp}, {0 + temp, 100 + temp},
		{100 + temp, 0 + temp}, {100 + temp, 100 + temp}})
	mt := gocv.GetPerspectiveTransform(pointVectorBefore, pointVectorAfter)
	projectionDst := gocv.NewMat()
	gocv.WarpPerspectiveWithParams(src, &projectionDst, mt, image.Point{200, 200}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})
	fmt.Println(projectionDst.Size())

	return projectionDst
}

// opencv 自带的qrcode检测识别，其对图片的要求较高，因此识别精准度不高。
// 要求二维码是正正方方的，不能旋转，不能投影
// 即便是干净的方正的二维码，草料可以识别，opencv是解析不了
//
// 于是增加了 github.com/makiuchi-d/gozxing ，此包实现了 zxing 的算法，效果较好，在 opencv都解析不出的时候它能解析出来
func DetectAndDecodeQrcode(input gocv.Mat) (content string, err error) {
	qrcodeDetector := gocv.NewQRCodeDetector()
	points := gocv.NewMat()
	exist := qrcodeDetector.Detect(input, &points)

	if exist {
		// log.Println(points.Size(), points.Type()) // [1 4] CV32FC2
		// points 为包围二维码的最小方框的四个点的坐标，依次为：左上，右上，右下，左下
		// data, _ := points.DataPtrFloat32()
		// log.Println(data) // [15 19 96 19 96 100 15 100]

		// reg := input.Region(image.Rect(int(data[0]), int(data[1]), int(data[4]), int(data[5])))
		// util.ShowImage("reg", reg, false)

		dst := gocv.NewMat()
		// dst 就是上面的 reg ，就是截取后的二维码图片，二值化后的。
		// content 二维码内容
		content = qrcodeDetector.Decode(input, points, &dst)
		// util.ShowImage("DetectAndDecodeQrcode", dst, true)
	} else {
		err = errors.New("qrcode not exist")
	}

	return
}

// opencv 中的 qrcode 与 gozxing 中的 qrcode 比较。
// opencv 自带的qrcode检测识别，其对图片的要求较高，因此识别精准度不高。
// 要求二维码是正正方方的，不能旋转，不能投影。
// 即便是干净的方正的二维码，opencv有时候解析不了。
// 于是增加了 github.com/makiuchi-d/gozxing ，此包实现了 zxing 的算法，效果较好，在 opencv都解析不出的时候它能解析出来。
// gozxing 会自动识别旋转，透视投影，缩放等行为，可以视为一个完整的解决方案。
func Run2() {
	input := gocv.IMRead("qrcoderecognition/4.jpg", gocv.IMReadColor)
	content, err := DetectAndDecodeQrcode(input)
	fmt.Println(content, err)

	img, err := input.ToImage()
	if err != nil {
		log.Fatal(err)
	}
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		log.Fatal(err)
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
