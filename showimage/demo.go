package showimage

import (
	"fmt"
	"image"
	"log"

	"github.com/phprao/go-opencv/util"
	"gocv.io/x/gocv"
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
	elem := gocv.GetStructuringElement(gocv.MorphRect, image.Point{15, 15})

	dst := gocv.NewMat()
	gocv.Erode(img, &dst, elem) // 腐蚀操作
	util.ShowImage("图像腐蚀-后", dst, false)

	dst2 := gocv.NewMat()
	gocv.Dilate(img, &dst2, elem) // 膨胀操作
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

	imgs := gocv.Split(srcImage) // B G R
	// imgSize := srcImage.Size()

	util.ShowImage("mat split", imgs[1], true)

	// fmt.Println(imgs[0].Size(), imgs[0].Channels(), imgs[0].Total())
	util.ShowMultipleImage("mat split", []gocv.Mat{imgs[0], imgs[1], imgs[2]}, 2)

	// befor := make([]uint8, srcImage.Total())

	// bm, _ := imgs[0].DataPtrUint8()
	// gm, _ := imgs[1].DataPtrUint8()
	// rm, _ := imgs[2].DataPtrUint8()

	// bm = append(append(bm, befor...), befor...)
	// gm = append(append(befor, gm...), befor...)
	// rm = append(append(befor, befor...), rm...)

	// b, _ := gocv.NewMatWithSizesFromBytes(imgSize, srcImage.Type(), bm)
	// g, _ := gocv.NewMatWithSizesFromBytes(imgSize, srcImage.Type(), gm)
	// r, _ := gocv.NewMatWithSizesFromBytes(imgSize, srcImage.Type(), rm)

	// fmt.Println(b.Size(), b.Channels(), b.Total())
	// fmt.Println(g.Size(), g.Channels(), g.Total())
	// fmt.Println(r.Size(), r.Channels(), r.Total())

	// util.ShowMultipleImage("mat split", []gocv.Mat{b, g, r}, 2)
}
