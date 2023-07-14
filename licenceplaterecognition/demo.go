package licenceplaterecognition

// 车牌识别

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"sort"

	"github.com/phprao/go-opencv/util"
	"gocv.io/x/gocv"
)

func Run() {
	src := gocv.IMRead("licenceplaterecognition/car_band.jpeg", gocv.IMReadColor)

	// 规范大小
	resize := gocv.NewMat()
	gocv.Resize(src, &resize, image.Point{620, 480}, 0, 0, gocv.InterpolationLinear)

	gray := gocv.NewMat()
	gocv.CvtColor(resize, &gray, gocv.ColorBGRToGray)

	// 使用双边滤波过滤掉不需要的细节
	filtered := gocv.NewMat()
	gocv.BilateralFilter(gray, &filtered, 13, 15, 15)

	// 边缘检测
	cannyed := gocv.NewMat()
	// 仅显示强度梯度大于最小阈值且小于最大阈值的边缘
	gocv.Canny(filtered, &cannyed, 30, 200)

	// 寻找轮廓
	originPointsVector := gocv.FindContours(cannyed, gocv.RetrievalTree, gocv.ChainApproxSimple)
	originPointsVectorSize := originPointsVector.Size()
	// fmt.Println("originPointsVectorSize", originPointsVectorSize) // 87

	// 一个二维数组，包含有87个一维数组，每一个数组下面包含了多个点，这些点组成当前轮廓。也就是有87个轮廓。
	// [[(10,12),(100,101),(23,56)...],[...],[...]]
	// fmt.Println(originPointsVector.ToPoints())

	// 把轮廓画出来，-1 表示画出所有轮廓
	drawBoundaryOnGray := gocv.NewMatWithSize(gray.Rows(), gray.Cols(), gocv.MatTypeCV8UC3)
	gocv.DrawContours(&drawBoundaryOnGray, originPointsVector, -1, color.RGBA{255, 0, 0, 255}, 1)
	// util.ShowImage("drawBoundaryOnGray", drawBoundaryOnGray, true)

	// 求得每一个轮廓的面积
	areas := make([]float64, 0)
	areasOld := make([]float64, 0)
	for i := 0; i < originPointsVectorSize; i++ {
		a := gocv.ContourArea(originPointsVector.At(i))
		areas = append(areas, a)
		areasOld = append(areasOld, a)
	}
	// 按面积排序，从小到大
	sort.Float64s(areas)
	// 提取出面积最大的10个轮廓
	areas = areas[originPointsVectorSize-10:]
	newPointsVector := gocv.NewPointsVector()
	for _, v := range areas {
		for k1, v1 := range areasOld {
			if v == v1 {
				newPointsVector.Append(originPointsVector.At(k1))
			}
		}
	}
	fmt.Println("newPointsVectorSize", newPointsVector.Size()) // 11，因为有两个面积相同
	drawNewPointsVector := gocv.NewMatWithSize(gray.Rows(), gray.Cols(), gocv.MatTypeCV8UC3)
	gocv.DrawContours(&drawNewPointsVector, newPointsVector, -1, color.RGBA{255, 0, 0, 255}, 1)
	// util.ShowImage("drawNewPointsVector", drawNewPointsVector, false)

	bandPointsVector := gocv.NewPointsVector()
	// 寻找具有四个表面闭合的轮廓
	for i := 0; i < newPointsVector.Size(); i++ {
		vetor := newPointsVector.At(i)
		peri := gocv.ArcLength(vetor, true)
		approx := gocv.ApproxPolyDP(vetor, 0.018*peri, true)
		if approx.Size() == 4 {
			bandPointsVector.Append(approx)
			break
		}
	}
	if bandPointsVector.IsNil() {
		log.Fatal("car band not found")
	}

	// 在原图上画个框
	gocv.DrawContours(&resize, bandPointsVector, -1, color.RGBA{255, 0, 0, 255}, 2)
	// util.ShowImage("drawBandPointsVectorOnResize", resize, false)

	// 抠出车牌
	mask := gocv.Zeros(gray.Rows(), gray.Cols(), gocv.MatTypeCV8U)
	// thickness = -1 表示填充
	gocv.DrawContours(&mask, bandPointsVector, -1, color.RGBA{255, 255, 255, 255}, -1)
	// util.ShowImage("mask", mask, false)
	newResize := gocv.NewMat()
	// gocv.BitwiseAndWithMask(resize, resize, &newResize, mask)
	resize.CopyToWithMask(&newResize, mask)
	// util.ShowImage("newResize", newResize, true)

	// 构造 ROI，先找到 minX, minY, maxX, maxY
	bandPoints := bandPointsVector.ToPoints()
	// fmt.Println(bandPoints) // [[(567,272) (556,427) (102,432) (99,277)]]
	var minX, maxX, minY, maxY int
	ys := make([]int, 4)
	for k, point := range bandPoints[0] {
		if point.X > maxX {
			maxX = point.X
		}
		if minX == 0 {
			minX = point.X
		} else if point.X < minX {
			minX = point.X
		}
		if point.Y > maxY {
			maxY = point.Y
		}
		if minY == 0 {
			minY = point.Y
		} else if point.Y < minY {
			minY = point.Y
		}
		ys[k] = point.Y
	}
	// fmt.Println(minX, maxX, minY, maxY) // 99 567 272 432

	band := gray.Region(image.Rect(minX, minY, maxX, maxY))
	gocv.Resize(band, &band, image.Point{460, 160}, 0, 0, gocv.InterpolationLinear)
	// util.ShowImage("band", band, false)

	// 车牌校正：车牌存在水平倾斜、垂直倾斜或梯形畸变等变形
	// 有了 bandPoints 之后就能根据三角函数算出倾斜的角度，然后沿着 band 中心点旋转来校正
	// 找到最上边的线与X轴的夹角
	sort.Ints(ys)
	var p0, p1 image.Point
	for _, point := range bandPoints[0] {
		if point.Y == ys[0] {
			p0 = point
		}
		if point.Y == ys[1] {
			p1 = point
		}
	}
	if p1.X < p0.X {
		p0, p1 = p1, p0
	}
	dy := math.Abs(float64(p1.Y - p0.Y))
	length := math.Sqrt(math.Pow(float64(p1.X-p0.X), 2) + math.Pow(float64(p1.Y-p0.Y), 2))
	radians := math.Acos(dy / length) // 范围：0 到 π
	degree := 180 * (math.Pi/2 - radians) / math.Pi
	if p1.Y < p0.Y {
		degree = 360 - degree
	}
	// fmt.Println(p0, p1, degree)
	mt := gocv.GetRotationMatrix2D(image.Point{band.Cols() / 2, band.Rows() / 2}, degree, 1) // 以图像中心点旋转
	gocv.WarpAffineWithParams(band, &band, mt, image.Point{band.Cols(), band.Rows()}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{0, 0, 0, 0})

	// util.ShowImage("bandrotate", band, true)

	// 二值化
	bandBinary := gocv.NewMat()
	gocv.Threshold(band, &bandBinary, 100, 255, gocv.ThresholdBinary)
	// util.ShowImage("bandBinary", bandBinary, false)
	bandBlur := gocv.NewMat()
	// 腐蚀
	gocv.Erode(bandBinary, &bandBlur, gocv.GetStructuringElement(gocv.MorphRect, image.Point{5, 5}))
	// util.ShowImage("bandBlur", bandBlur, false)
	// 模糊
	gocv.Blur(bandBlur, &bandBlur, image.Point{5, 5})
	// util.ShowImage("bandBlur", bandBlur, false)
	gocv.Dilate(bandBlur, &bandBlur, gocv.GetStructuringElement(gocv.MorphRect, image.Point{5, 5}))
	util.ShowImage("bandBlur", bandBlur, true)

	// 字符分割：垂直投影法
	colsMap := make([]int, bandBlur.Cols())
	for r := 0; r < bandBlur.Rows(); r++ {
		for c := 0; c < bandBlur.Cols(); c++ {
			if bandBlur.GetUCharAt(r, c) > 0 {
				colsMap[c]++
			}
		}
	}
	// fmt.Println(colsMap)
	// 画个直方图
	util.SaveHistSingle[int](colsMap, "垂直投影", "边界", "licenceplaterecognition/band.html")

	// pytesseract是基于Python的OCR工具， 底层使用的是Google的Tesseract-OCR 引擎，支持识别图片中的文字，支持jpeg, png, gif, bmp, tiff等图片格式。
	// text = pytesseract.image_to_string(band, config='--psm 11')

}
