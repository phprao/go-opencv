package showimage

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"

	"github.com/phprao/go-opencv/util"
	"gocv.io/x/gocv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func Run() {
	w := gocv.NewWindow("show image")
	w.MoveWindow(400, 300)
	util.ReadAndShowImage(w, "showimage/cat.jpg")

	log.Println(w.WaitKey(10000))
}

// 图像腐蚀与膨胀
func Run3() {
	img := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)
	util.ShowImage("原图", img, false)

	// 设置腐蚀块大小，gocv中使用 image.Point 来设置宽高
	// 返回 15*15 的矩阵，元素值都是1，用于形态学（morphologic）的操作
	elem := gocv.GetStructuringElement(gocv.MorphRect, image.Point{15, 15})

	dst := gocv.NewMat()
	gocv.Erode(img, &dst, elem) // 腐蚀操作，黑色变大
	util.ShowImage("图像腐蚀-后", dst, false)

	dst2 := gocv.NewMat()
	gocv.Dilate(img, &dst2, elem) // 膨胀操作，白色变大
	util.ShowImage("图像膨胀-后", dst2, true)
}

// 边缘检测
func Run4() {
	srcImage := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)

	// 将原始图像转换为灰度图像
	grayImage := gocv.NewMat()
	gocv.CvtColor(srcImage, &grayImage, gocv.ColorBGRToGray) // 单通道

	// 先用3*3内核来降噪，模糊处理
	edge := gocv.NewMat()
	gocv.Blur(grayImage, &edge, image.Point{3, 3})

	// 运行canny算子
	dstImage := gocv.NewMat()
	gocv.Canny(edge, &dstImage, 3, 9)

	util.ShowMultipleImage("边缘检测", []gocv.Mat{srcImage, grayImage, edge, dstImage}, 2)
}

// 卷积运算
func Run4_1() {
	data := []byte{181, 8, 127, 14, 208, 158, 144, 59, 51, 179, 228, 118, 160, 7, 212, 101, 242, 240, 151, 136, 22, 95, 63, 47, 250, 215, 161, 189, 66, 234, 146, 166, 70, 201, 136, 201}
	src, _ := gocv.NewMatWithSizesFromBytes([]int{6, 6}, gocv.MatTypeCV8U, data)
	fmt.Println(src.DataPtrUint8()) // 6*6

	kernel := gocv.NewMatWithSize(3, 3, gocv.MatTypeCV32F)

	kernel.SetFloatAt(0, 0, 0.1111)
	kernel.SetFloatAt(0, 2, 0.1111)
	kernel.SetFloatAt(1, 0, 0.1111)
	kernel.SetFloatAt(1, 2, 0.1111)
	kernel.SetFloatAt(2, 0, 0.1111)
	kernel.SetFloatAt(2, 2, 0.1111)
	fmt.Println(kernel.DataPtrFloat32()) // 3*3

	dst := gocv.NewMat()
	gocv.Filter2D(src, &dst, gocv.MatTypeCV32F, kernel, image.Point{-1, -1}, 0, gocv.BorderConstant)
	fmt.Println(dst.DataPtrFloat32()) // 6*6
	fmt.Println(dst.Size())
}

// 扩充边界
func Run4_2() {
	data := []byte{181, 8, 127, 14, 208, 158, 144, 59, 51, 179, 228, 118, 160, 7, 212, 101, 242, 240, 151, 136, 22, 95, 63, 47, 250, 215, 161, 189, 66, 234, 146, 166, 70, 201, 136, 201}
	src, _ := gocv.NewMatWithSizesFromBytes([]int{6, 6}, gocv.MatTypeCV8U, data)
	fmt.Println(src.Size()) // 6*6

	srcWithBorder := gocv.NewMat()
	gocv.CopyMakeBorder(src, &srcWithBorder, 1, 1, 1, 1, gocv.BorderConstant, color.RGBA{0, 0, 0, 0})
	fmt.Println(srcWithBorder.Size()) // 8*8
}

// 图像翻转
func Run5() {
	srcImage := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)
	util.ShowImage("翻转-前", srcImage, false)

	dstImage := gocv.NewMat()
	// 0 - 沿着水平线翻转
	// 1 - 沿着垂直线翻转
	// -1 - 沿着水平和垂直线翻转
	gocv.Flip(srcImage, &dstImage, -1)

	util.ShowImage("翻转-后", dstImage, true)
}

// 图像阈值化
func Run6() {
	srcImage := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)
	util.ShowImage("原图", srcImage, false)

	grayImage := gocv.NewMat()
	gocv.CvtColor(srcImage, &grayImage, gocv.ColorBGRToGray)
	util.ShowImage("灰度", grayImage, false)

	dstImage := gocv.NewMat()
	gocv.Threshold(grayImage, &dstImage, 125, 255, gocv.ThresholdBinary)
	util.ShowImage("ThresholdBinary", dstImage, true)
}

// 关于行数，列数，通道的关系
// Total 是像素点的个数，row * col
// 矩阵变换
func Run7() {
	m1 := gocv.NewMatWithSize(20, 30, gocv.MatTypeCV8UC1)
	fmt.Printf("size:%v, elemSize:%v, type:%v, total:%v, channels:%v\n", m1.Size(), m1.ElemSize(), m1.Type(), m1.Total(), m1.Channels())
	// 	size:[20 30], elemSize:1, type:CV8U, total:600, channels:1

	// cn int 通道数，0表示保持原通道数不变
	// rows int 矩阵行数，0表示保持原行数不变，列数会自动计算
	// 变换规则：row1 * col1 * channel1 == row2 * col2 * channel2
	m2 := m1.Reshape(2, 20) // 2通道，20行N列
	fmt.Printf("size:%v, elemSize:%v, type:%v, total:%v, channels:%v\n", m2.Size(), m2.ElemSize(), m2.Type(), m2.Total(), m2.Channels())
	// size:[20 15], elemSize:2, type:CV8UC2, total:300, channels:2

	m3 := m1.Reshape(1, 1) // 1通道，1行N列
	fmt.Printf("size:%v, elemSize:%v, type:%v, total:%v, channels:%v\n", m3.Size(), m3.ElemSize(), m3.Type(), m3.Total(), m3.Channels())
	// size:[1 600], elemSize:1, type:CV8U, total:600, channels:1

	m4 := m3.T() // 转置操作，得到 N行1列，1通道
	fmt.Printf("size:%v, elemSize:%v, type:%v, total:%v, channels:%v\n", m4.Size(), m4.ElemSize(), m4.Type(), m4.Total(), m4.Channels())
	// size:[600 1], elemSize:1, type:CV8U, total:600, channels:1

	m5, err := gocv.NewMatFromBytes(2, 3, gocv.MatTypeCV8UC3, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("size:%v, elemSize:%v, type:%v, total:%v, channels:%v\n", m5.Size(), m5.ElemSize(), m5.Type(), m5.Total(), m5.Channels())
	// size:[2 3], elemSize:3, type:CV8UC3, total:6, channels:3

	m6 := m5.RowRange(0, 1) // 获取部分行组成新矩阵
	fmt.Printf("size:%v, elemSize:%v, type:%v, total:%v, channels:%v\n", m6.Size(), m6.ElemSize(), m6.Type(), m6.Total(), m6.Channels())
	// size:[1 3], elemSize:3, type:CV8UC3, total:3, channels:3

	// 打印矩阵数据，元素的个数为 row * col * channel
	sli, err := m6.DataPtrUint8()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sli) // [1 2 3 4 5 6 7 8 9]

	fmt.Println(m1.Step()) // 30，返回每一行占用的字节数
}

// 读物GIF动图
func Run8() {
	util.ReadAndShowGIF("showimage/image15.gif")
}

// 读取视频文件
func Run2() {
	util.ReadAndShowVideo("showimage/video1.mp4")
}

// 读取网络图片
func Run9() {
	util.ReadAndShowImageFromUrl("https://videoactivity.bookan.com.cn/ac_1_1687763619_793.jpg")
}

// 图片裁剪，使用鼠标选择感兴趣的区域
func Run10() {
	w := gocv.NewWindow("image")
	srcImage := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)
	rect := w.SelectROI(srcImage)
	subImage := srcImage.Region(rect)
	w.IMShow(subImage)
	w.WaitKey(0)
}

// 拆分通道
func Run11() {
	srcImage := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)

	// 拆分通道，此处将3通道的数据拆分成，三个单通道的数据，依次为 B G R
	imgs := gocv.Split(srcImage)

	bm := imgs[0].ToBytes()
	fmt.Println("bm", bm[:12])
	// bm [53 54 54 55 55 56 56 56 58 59 62 65]

	gm := imgs[1].ToBytes()
	fmt.Println("gm", gm[:12])
	// gm [55 56 56 57 57 58 58 58 59 60 63 66]

	rm := imgs[2].ToBytes()
	fmt.Println("rm", rm[:12])
	// rm [95 96 96 97 98 99 99 99 103 104 107 110]

	sm := srcImage.ToBytes()
	fmt.Println("sm", sm[:36])
	// sm [53 55 95 54 56 96 54 56 96 55 57 97 55 57 98 56 58 99 56 58 99 56 58 99 58 59 103 59 60 104 62 63 107 65 66 110]
}

// 拆分通道
func Run12() {
	srcImage := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)
	imgs := gocv.Split(srcImage)

	basicb := make([]byte, srcImage.Channels()*srcImage.Total())
	bm := imgs[0].ToBytes()
	for k := range basicb {
		if k%3 == 0 {
			basicb[k] = bm[k/3]
		}
	}
	b, err := gocv.NewMatWithSizesFromBytes(srcImage.Size(), srcImage.Type(), basicb)
	if err != nil {
		log.Fatal(err)
	}

	basicg := make([]byte, srcImage.Channels()*srcImage.Total())
	gm := imgs[1].ToBytes()
	for k := range basicg {
		if k%3 == 1 {
			basicg[k] = gm[k/3]
		}
	}
	g, err := gocv.NewMatWithSizesFromBytes(srcImage.Size(), srcImage.Type(), basicg)
	if err != nil {
		log.Fatal(err)
	}

	basicr := make([]byte, srcImage.Channels()*srcImage.Total())
	rm := imgs[2].ToBytes()
	for k := range basicr {
		if k%3 == 2 {
			basicr[k] = rm[k/3]
		}
	}
	r, err := gocv.NewMatWithSizesFromBytes(srcImage.Size(), srcImage.Type(), basicr)
	if err != nil {
		log.Fatal(err)
	}

	util.ShowMultipleImage("mat split", []gocv.Mat{srcImage, b, g, r}, 2)
}

// 拆分通道
func Run13() {
	srcImage := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)
	imgs := gocv.Split(srcImage) // B G R

	basic := make([]byte, srcImage.Total())
	basicMat, err := gocv.NewMatWithSizesFromBytes(srcImage.Size(), gocv.MatTypeCV8UC1, basic)
	if err != nil {
		log.Fatal(err)
	}

	b := gocv.NewMat()
	gocv.Merge([]gocv.Mat{imgs[0], basicMat, basicMat}, &b)
	g := gocv.NewMat()
	gocv.Merge([]gocv.Mat{basicMat, imgs[1], basicMat}, &g)
	r := gocv.NewMat()
	gocv.Merge([]gocv.Mat{basicMat, basicMat, imgs[2]}, &r)

	util.ShowMultipleImage("mat split", []gocv.Mat{srcImage, b, g, r}, 2)
}

// Scalar
func Run14() {
	scalar := gocv.NewScalar(255, 0, 0, 0)
	m := gocv.NewMatWithSizesWithScalar([]int{200, 200}, gocv.MatTypeCV8UC3, scalar)
	fmt.Println(m.Size(), m.Channels(), m.Total(), m.ElemSize())
	util.ShowImage("Scalar", m, true)
}

// 使用遮罩
func Run15() {
	// mask := gocv.NewMat()
	// mask.SetUCharAt(0, 0, 255)

	// CopyToWithMask
}

// 垂直拼接，水平拼接
func Run16() {
	scalar1 := gocv.NewScalar(255, 0, 0, 0)
	m1 := gocv.NewMatWithSizesWithScalar([]int{200, 200}, gocv.MatTypeCV8UC3, scalar1)

	scalar2 := gocv.NewScalar(0, 255, 0, 0)
	m2 := gocv.NewMatWithSizesWithScalar([]int{200, 200}, gocv.MatTypeCV8UC3, scalar2)

	scalar3 := gocv.NewScalar(0, 0, 255, 0)
	m3 := gocv.NewMatWithSizesWithScalar([]int{200, 200}, gocv.MatTypeCV8UC3, scalar3)

	dst1 := gocv.NewMat()
	// 沿水平堆叠，要求行数必须一致
	gocv.Hconcat(m1, m2, &dst1)
	util.ShowImage("Hconcat", dst1, false)

	dst2 := gocv.NewMat()
	// 垂直方向堆叠，要求列数必须相同，因此此处会报错
	gocv.Vconcat(dst1, m3, &dst2)

	util.ShowImage("Vconcat", dst2, true)
}

// 图像加法，不同尺寸相加
func Run17() {
	iconImage := gocv.IMRead("showimage/gongfu.jpg", gocv.IMReadColor)
	iconSize := iconImage.Size()
	backImage := gocv.IMRead("showimage/girl.jpg", gocv.IMReadColor)

	alpha, beta, gamma := 1.0, 0.0, 0.0
	backRegin := backImage.Region(image.Rect(10, 10, 10+iconSize[1], 10+iconSize[0]))
	dstImage := gocv.NewMat()

	// add, addWeight 要求尺寸和通道数完全一致
	// alpha, beta 分别表示第一个和第二个的权重，0-1之间，通常 alpha+beta=1
	// gamma 灰度系数，图像校正的偏移量，用于调节亮度，dst = src1 * alpha + src2 * beta + gamma
	gocv.AddWeighted(iconImage, alpha, backRegin, beta, gamma, &dstImage)

	dstImage.CopyTo(&backRegin)
	util.ShowImage("addWeight", backImage, true)
}

// 圆形蒙版
func Run18() {
	src := gocv.IMRead("showimage/girl.jpg", gocv.IMReadColor)
	// 蒙版尺寸与源图像的尺寸要一致
	// 蒙版为 8位灰度格式，0表示过滤掉，255表示接受
	mask := gocv.NewMatWithSize(src.Size()[0], src.Size()[1], gocv.MatTypeCV8U)
	regin := mask.Region(image.Rect(100, 100, 200, 200))
	regin.AddUChar(255)
	dst := gocv.NewMat()
	src.CopyToWithMask(&dst, mask)

	// 圆形蒙版
	mask2 := gocv.NewMatWithSize(src.Size()[0], src.Size()[1], gocv.MatTypeCV8U)
	// thickness 为线条粗细，-1表示填充
	gocv.Circle(&mask2, image.Point{250, 150}, 50, color.RGBA{255, 255, 255, 255}, -1)
	src.CopyToWithMask(&dst, mask2)

	util.ShowImage("CopyToWithMask", dst, true)
}

// 把 jpg 中的图标扣出来粘贴到另一张图片上
func Run19() {
	// logo
	logoSrc := gocv.IMRead("showimage/logo.jpg", gocv.IMReadColor)

	// 灰度化
	logoGray := gocv.NewMat()
	gocv.CvtColor(logoSrc, &logoGray, gocv.ColorBGRToGray)

	// 二值化
	logoBin := gocv.NewMat()
	gocv.Threshold(logoGray, &logoBin, 175, 255, gocv.ThresholdBinary)

	// 黑白颠倒
	logoInv := gocv.NewMat()
	gocv.BitwiseNot(logoBin, &logoInv)

	// 抠图
	logo := gocv.NewMat()
	logoSrc.CopyToWithMask(&logo, logoInv)

	// 目标
	PageSrc := gocv.IMRead("showimage/girl.jpg", gocv.IMReadColor)
	regin := PageSrc.Region(image.Rect(0, 0, logoSrc.Cols(), logoSrc.Rows()))

	regin2 := gocv.NewMat()
	regin.CopyToWithMask(&regin2, logoBin)

	regin3 := gocv.NewMat()
	gocv.Add(logo, regin2, &regin3)

	util.ShowMultipleImage("logo", []gocv.Mat{logoSrc, logoGray, logoBin, logoInv, logo, regin2, regin3}, 3)

	regin3.CopyTo(&regin)

	util.ShowImage("logo2", PageSrc, true)
}

// 图像中添加文本
func Run20() {
	PageSrc := gocv.IMRead("showimage/girl.jpg", gocv.IMReadColor)
	// fontFace 字体
	// fontScale 字体缩放比例
	// tnickness 线条宽度
	// bottomLeftOrigin true表示数据原点位于左下角，False表示位于左上角，true是上下颠倒的
	// 不支持中文（包括中文标点符号）
	gocv.PutTextWithParams(&PageSrc, "OpenCV 2023, showimage/girl.jpg", image.Point{100, 100}, gocv.FontHersheySimplex, 1, color.RGBA{0, 255, 0, 0}, 3, gocv.Filled, false)

	// 关于中文文本，需要使用 github.com/golang/freetype
	// 此处的 fontSize 要设置大一些，否则看不见
	err := util.WriteTextOnMat(&PageSrc, "关于中文文本，需要使用freetype", image.Point{100, 300}, "c:/windows/fonts/msyh.ttc", 30, color.RGBA{0, 255, 0, 255})
	if err != nil {
		fmt.Println(err)
	}

	util.ShowImage("show text", PageSrc, true)
}

// 图像平移
func RUn21() {
	// 二维空间的变换矩阵为 2行3列

	// 平移
	// vec1 = M·vec0 ； M 为变换矩阵，跟 OpenGL 一样，dx为正则向右，dy为正则向下
	/*
			| 1 0 dx |			| x	|			 | x + dx|
		M = | 0 1 dy |     vec =| y |      vec1 =| y + dy|
			| 0 0 1  |			| 1 |			 | 1     |
	*/
	img := gocv.IMRead("showimage/gongfu.jpg", gocv.IMReadColor)

	imgDst := gocv.NewMat()
	mt := gocv.NewMatWithSize(2, 3, gocv.MatTypeCV32F)
	mt.SetFloatAt(0, 0, 1) // row, col 都是从0开始的
	mt.SetFloatAt(0, 2, 10)
	mt.SetFloatAt(1, 1, 1)
	mt.SetFloatAt(1, 2, 10)
	fmt.Println(mt.DataPtrFloat32())

	// m 变换矩阵，要求 float32 类型，否则会报错
	// sz 输出的图像尺寸
	// flags 图像的插值算法，默认为 gocv.InterpolationLinear
	// borderType 边界像素的计算方法，默认为 gocv.BorderConstant
	// borderValue 边界填充值
	gocv.WarpAffineWithParams(img, &imgDst, mt, image.Point{img.Cols(), img.Rows()}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})

	util.ShowMultipleImage("translate image", []gocv.Mat{img, imgDst}, 2)
}

// 图像旋转
func Run22() {
	img := gocv.IMRead("showimage/gongfu.jpg", gocv.IMReadColor)
	img90 := gocv.NewMat()
	img180 := gocv.NewMat()
	img270 := gocv.NewMat()
	// 围绕着图像中心点 90,180,270
	gocv.Rotate(img, &img90, gocv.Rotate90Clockwise)
	gocv.Rotate(img, &img180, gocv.Rotate180Clockwise)
	gocv.Rotate(img, &img270, gocv.Rotate90CounterClockwise)

	util.ShowMultipleImage("translate image", []gocv.Mat{img, img90, img180, img270}, 3)
}

// 图像旋转
func Run23() {
	// 围绕着左上角（0,0）点旋转
	/*
			| cosθ  -sinθ  0 |
		M = | sinθ  cosθ   0 |
			| 0     0      1 |
	*/
	img := gocv.IMRead("showimage/gongfu.jpg", gocv.IMReadColor)
	θ := math.Pi / 6
	imgDst := gocv.NewMat()
	mt := gocv.NewMatWithSize(2, 3, gocv.MatTypeCV32F)
	mt.SetFloatAt(0, 0, float32(math.Cos(θ))) // row, col 都是从0开始的
	mt.SetFloatAt(0, 1, -float32(math.Sin(θ)))
	mt.SetFloatAt(1, 0, float32(math.Sin(θ)))
	mt.SetFloatAt(1, 1, float32(math.Cos(θ)))

	gocv.WarpAffineWithParams(img, &imgDst, mt, image.Point{img.Cols(), img.Rows()}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})

	util.ShowMultipleImage("translate image", []gocv.Mat{img, imgDst}, 2)
}

// 图像旋转
func Run24() {
	// 围绕着任意点旋转，先将该点移动到（0,0）处，旋转，再将该点反向移回
	/*
			| 1 0 dx |   | cosθ  -sinθ  0 |   | 1 0 -dx |
		M = | 0 1 dy | · | sinθ  cosθ   0 | · | 0 1 -dy |
			| 0 0 1  |   | 0     0      1 |   | 0 0  1  |
	*/
	img := gocv.IMRead("showimage/gongfu.jpg", gocv.IMReadColor)
	var θ float64 = 45
	imgDst := gocv.NewMat()

	// scale 为缩放比例
	// angle 为度数，而不是弧度
	// 注意是逆时针旋转
	mt := gocv.GetRotationMatrix2D(image.Point{img.Cols() / 2, img.Rows() / 2}, θ, 1) // 以图像中心点旋转

	gocv.WarpAffineWithParams(img, &imgDst, mt, image.Point{img.Cols(), img.Rows()}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})

	util.ShowMultipleImage("translate image", []gocv.Mat{img, imgDst}, 2)
}

// 图像翻转
func Run25() {
	img := gocv.IMRead("showimage/gongfu.jpg", gocv.IMReadColor)
	imgDst1 := gocv.NewMat()
	imgDst2 := gocv.NewMat()
	imgDst3 := gocv.NewMat()

	// horizontal(0), vertical(1), or both axes(-1)
	gocv.Flip(img, &imgDst1, 0)
	gocv.Flip(img, &imgDst2, 1)
	gocv.Flip(img, &imgDst3, -1)

	util.ShowMultipleImage("translate image", []gocv.Mat{img, imgDst1, imgDst2, imgDst3}, 2)
}

// 图像金字塔
func Run26() {
	img := gocv.IMRead("showimage/girl.jpg", gocv.IMReadColor) // 800 * 732

	imgDst := make([]gocv.Mat, 0)
	imgCopy := gocv.NewMat()
	img.CopyTo(&imgCopy)

	for i := 1; i <= 3; i++ {
		dst := gocv.NewMat()
		gocv.PyrDown(imgCopy, &dst, image.Point{imgCopy.Rows() / 2, imgCopy.Cols() / 2}, gocv.BorderDefault)
		imgDst = append(imgDst, dst)
		fmt.Println(dst.Size())
		dst.CopyTo(&imgCopy)
	}

	util.ShowImage("imgDst1", imgDst[0], false)
	util.ShowImage("imgDst2", imgDst[1], false)
	util.ShowImage("imgDst3", imgDst[2], true)
}

// 错切，斜切，扭变
func Run27() {
	img := gocv.IMRead("showimage/box.jpg", gocv.IMReadColor)
	imgDst := gocv.NewMat()
	θ := math.Pi / 12
	mt := gocv.NewMatWithSize(2, 3, gocv.MatTypeCV32F)
	mt.SetFloatAt(0, 0, 1) // row, col 都是从0开始的
	mt.SetFloatAt(0, 1, float32(math.Tan(θ)))
	mt.SetFloatAt(1, 1, 1)

	gocv.WarpAffineWithParams(img, &imgDst, mt, image.Point{img.Cols(), img.Rows()}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})

	util.ShowMultipleImage("translate image", []gocv.Mat{img, imgDst}, 2)
}

// 投影变换
// Perspective 透视投影
func Run28() {
	img := gocv.IMRead("showimage/box.jpg", gocv.IMReadColor)
	imgW := img.Cols()
	imgH := img.Rows()
	imgDst := gocv.NewMat()

	// 根据图像中不共线的四个点在变换前后的对应位置求得3*3的变换矩阵
	pointVectorBefore := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {imgW, 0}, {imgW, imgH}, {0, imgH}})
	pointVectorAfter := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {imgW, 0}, {imgW - 50, imgH - 50}, {50, imgH - 50}})
	mt := gocv.GetPerspectiveTransform(pointVectorBefore, pointVectorAfter)

	// 透视投影
	gocv.WarpPerspectiveWithParams(img, &imgDst, mt, image.Point{img.Cols(), img.Rows()}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 255})

	util.ShowMultipleImage("Perspective image", []gocv.Mat{img, imgDst}, 2)
}

// 投影变换
func Run29() {
	img := gocv.IMRead("showimage/box.jpg", gocv.IMReadColor)
	imgW := img.Cols()
	imgH := img.Rows()
	imgDst := gocv.NewMat()

	// 根据图像中不共线的四个点在变换前后的对应位置求得3*3的变换矩阵
	pointVectorBefore := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {imgW, 0}, {imgW, imgH}, {0, imgH}})
	pointVectorAfter := gocv.NewPointVectorFromPoints([]image.Point{{0, 0}, {imgW / 2, 0}, {imgW/2 - 50, imgH/2 - 50}, {50, imgH/2 - 50}})
	mt := gocv.GetPerspectiveTransform(pointVectorBefore, pointVectorAfter)

	// 边界填充选项
	gocv.WarpPerspectiveWithParams(img, &imgDst, mt, image.Point{img.Cols(), img.Rows()}, gocv.InterpolationArea, gocv.BorderWrap, color.RGBA{255, 255, 255, 255})

	util.ShowMultipleImage("Perspective image", []gocv.Mat{img, imgDst}, 2)
}

// 直角坐标（笛卡尔坐标）与极坐标的转换
func Run30() {
	xd := []float32{0, 1, 2, 0, 1, 2, 0, 1, 2}
	yd := []float32{0, 0, 0, 1, 1, 1, 2, 2, 2}
	x := gocv.NewMatWithSize(1, 9, gocv.MatTypeCV32F)
	y := gocv.NewMatWithSize(1, 9, gocv.MatTypeCV32F)
	for k, v := range xd {
		x.SetFloatAt(0, k, v)
	}
	for k, v := range yd {
		y.SetFloatAt(0, k, v)
	}
	mag := gocv.NewMat()
	ang := gocv.NewMat()

	// x, y：直角坐标系的横坐标、纵坐标
	// magnitude, angle：极坐标系的向量值、角度值
	// angleInDegrees：弧度制/角度值选项，默认值 0 选择弧度制，1 选择角度制（[0,360]）
	// 直接坐标 --> 极坐标
	gocv.CartToPolar(x, y, &mag, &ang, true)

	fmt.Println(mag.DataPtrFloat32())
	fmt.Println(ang.DataPtrFloat32())

	//////////////////////////////////////////////////////////////////////////
	x1 := gocv.NewMat()
	y1 := gocv.NewMat()
	// 极坐标 --> 直接坐标
	gocv.PolarToCart(mag, ang, &x1, &y1, true)
	fmt.Println(x1.DataPtrFloat32())
	fmt.Println(y1.DataPtrFloat32())
}

// 直角坐标（笛卡尔坐标）与极坐标的转换
func Run31() {
	img := gocv.IMRead("showimage/circle.jpg", gocv.IMReadColor)
	dst := gocv.NewMat()

	cx := img.Cols() / 2
	cy := img.Rows() / 2
	mR := math.Max(float64(cx), float64(cy))
	gocv.LinearPolar(img, &dst, image.Point{cx, cy}, mR, gocv.InterpolationLinear)

	dstRotate := gocv.NewMat()
	gocv.Rotate(dst, &dstRotate, gocv.Rotate90CounterClockwise)

	util.ShowMultipleImage("LinearPolar", []gocv.Mat{img, dst, dstRotate}, 3)
}

// 比特平面分层
func Run32() {
	img := gocv.IMRead("showimage/cat.jpg", gocv.IMReadGrayScale)
	dst := make([]gocv.Mat, 9)
	dst[0] = img

	for i := 0; i < 8; i++ {
		d := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)

		for a := 0; a < img.Rows(); a++ {
			for b := 0; b < img.Cols(); b++ {
				grayByte := fmt.Sprintf("%08s", strconv.FormatInt(int64(img.GetUCharAt(a, b)), 2))
				it, _ := strconv.Atoi(string(grayByte[i]))
				it = it * 255 // 为了展示乘以255，实际存储的时候不需要

				d.SetUCharAt(a, b, uint8(it))
			}
		}
		dst[i+1] = d
	}

	util.ShowMultipleImage("bit-split", dst, 3)
}

// 灰度直方图 + gonum/plot
func Run33() {
	img := gocv.IMRead("showimage/cat.jpg", gocv.IMReadGrayScale)
	hist := gocv.NewMat()
	mask := gocv.NewMat()

	// src：输入图像
	// channels：直方图计算的通道
	// mask：蒙版
	// size：直方柱的数量，一般取 [256]
	// ranges：像素值的取值范围，一般为 [0,256]
	// acc: accumulate是否累加，一般为false
	gocv.CalcHist([]gocv.Mat{img}, []int{0}, mask, &hist, []int{256}, []float64{0, 256}, false)
	if hist.Empty() || hist.Rows() != 256 || hist.Cols() != 1 {
		log.Fatal("Invalid CalcHist test")
	}

	// fmt.Println(hist.Type()) ---> CV32F
	data, _ := hist.DataPtrFloat32()
	var max float32 = 0
	var kk int
	for k, v := range data {
		if v > max {
			max = v
			kk = k
		}
	}

	fmt.Println(kk, max) // 195 1348

	p := plot.New()

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	pts := make(plotter.XYs, 256)
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = float64(data[i])
	}

	err := plotutil.AddLinePoints(p, "CalcHist", pts)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "showimage/CalcHist.png"); err != nil {
		panic(err)
	}
}

// 灰度直方图 + go-echarts
func Run33_1() {
	img := gocv.IMRead("showimage/cat.jpg", gocv.IMReadGrayScale)
	hist := gocv.NewMat()
	mask := gocv.NewMat()

	// src：输入图像
	// channels：直方图计算的通道
	// mask：蒙版
	// size：直方柱的数量，一般取 [256]
	// ranges：像素值的取值范围，一般为 [0,256]
	// acc: accumulate是否累加，一般为false
	gocv.CalcHist([]gocv.Mat{img}, []int{0}, mask, &hist, []int{256}, []float64{0, 256}, false)
	if hist.Empty() || hist.Rows() != 256 || hist.Cols() != 1 {
		log.Fatal("Invalid CalcHist test")
	}

	// fmt.Println(hist.Type()) ---> CV32F
	data, _ := hist.DataPtrFloat32()

	// 直方图
	util.SaveHistSingle[float32](data, "灰度直方图", "灰度值", "showimage/Run33_1_bar.html")

	// 折线图
	// util.SaveLineSingle[float32](data, "灰度直方图", "灰度值", "showimage/Run33_1_bar.html", false)
}

// 直方图均衡化
func Run34() {
	img := gocv.IMRead("showimage/cat.jpg", gocv.IMReadGrayScale)
	dst := gocv.NewMat()
	gocv.EqualizeHist(img, &dst)

	util.ShowMultipleImage("EqualizeHist", []gocv.Mat{img, dst}, 2)
}

// 直方图均衡化
func Run34_1() {
	img := gocv.IMRead("showimage/cat.jpg", gocv.IMReadGrayScale)
	dst := gocv.NewMat()
	hist := gocv.NewMat()
	mask := gocv.NewMat()
	gocv.EqualizeHist(img, &dst)

	gocv.CalcHist([]gocv.Mat{dst}, []int{0}, mask, &hist, []int{256}, []float64{0, 256}, false)
	if hist.Empty() || hist.Rows() != 256 || hist.Cols() != 1 {
		log.Fatal("Invalid CalcHist test")
	}

	// fmt.Println(hist.Type()) ---> CV32F
	data, _ := hist.DataPtrFloat32()

	util.SaveHistSingle[float32](data, "灰度直方图", "灰度值", "showimage/Run34_1_bar.html")
}

// 反色变换
func Run35() {
	img := gocv.IMRead("showimage/cat.jpg", gocv.IMReadGrayScale)
	imgdata := img.ToBytes()
	l := len(imgdata)
	dst := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	for i := 0; i < l; i++ {
		dst.SetUCharAt(i/img.Cols(), i%img.Cols(), 255-imgdata[i])
	}

	util.ShowMultipleImage("Invert", []gocv.Mat{img, dst}, 2)
}
