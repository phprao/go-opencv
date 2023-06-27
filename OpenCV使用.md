`func NewWindow(name string) *Window`

创建窗口，可以创建多个，但是名称要不一样。

opencv中的GUI是由`HighGUI`这个组件提供的。

`func (w *Window) WaitKey(delay int) int`

用于阻塞当前线程，否则窗口会一闪而过。如果参数为0就会一直阻塞，直到有按键事件；如果大于0，就会等待相应的毫秒然后解除阻塞，或者有按键事件就会立即解除阻塞。返回值为按键的code，如果没有按键则返回值为-1。如果有多个窗口会同时生效。

`func IMRead(name string, flags IMReadFlag) Mat`

加载图片文件。常用的读取模式：1、`IMReadColor`读取出来的为彩色图片，拥有三个通道，opencv中的通道顺序为`BGR`，而不是`RGB`；2、`IMReadGrayScale`读取出来的为黑白照片，只有一个通道。支持的图片格式很多，比如常用的`jpg,jpeg,png,bmp,webp,awf`，具体可以看文档，不支持`gif`。

#### 案例

##### 图像腐蚀与膨胀

类似于马赛克的效果。

```go
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
```

![image-20230627170654181](D:\dev\php\magook\trunk\server\md\img\image-20230627170654181.png)

腐蚀的作用就是让暗的区域变大（它是在一个区域内找0），而膨胀的作用就是让亮的区域变大（它是在一个区域内找255）。而最终的结果很大程度上取决于“核”的形状和大小，“核”的形状有方形、X形、椭圆形。

“核”就是对比的对象，上例中，我们设置了15X15的方形区域，其实就是一个矩阵，该矩阵每个点填充1，然后拿着这个矩阵从图像左上角开始对比，符合一定的条件就将图像上的某个像素点置为0或者255。

最典型的一个应用场景就是在你二值化后，你的目标和背景抠的不是很干净，比如里面还有一些干扰点或者干扰线，可以试试两个操作，有时候效果出奇的好。



##### 边缘检测

通俗来讲就是试图勾勒出物体轮廓。

```go
func Run4() {
	srcImage := gocv.IMRead("showimage/cat.jpg", gocv.IMReadColor)
	util.ShowImage("边缘检测-前", srcImage, false)

	// 将原始图像转换为灰度图像
	grayImage := gocv.NewMat()
	gocv.CvtColor(srcImage, &grayImage, gocv.ColorBGRToGray)
	util.ShowImage("灰度", grayImage, false)

	// 先用3*3内核来降噪，模糊处理
	edge := gocv.NewMat()
	gocv.Blur(grayImage, &edge, image.Point{3, 3})
	util.ShowImage("降噪", edge, false)

	// 运行canny算子
	dstImage := gocv.NewMat()
	gocv.Canny(edge, &dstImage, 3, 9)

	util.ShowImage("边缘检测-后", dstImage, true)
}
```

![image-20230627154904325](D:\dev\php\magook\trunk\server\md\img\image-20230627154904325.png)

关于图像的`模糊处理/平滑处理/降噪处理/滤波处理`

在opencv2中，可能使用blur对图像进行平滑处理，这种方法就是最简单的求平均数。*平滑* 也称 *模糊*, 是一项简单且使用频率很高的图像处理方法。平滑处理的用途有很多， 但是在很多地方我们仅仅关注它减少噪声的功用。平滑处理时需要用到一个 *滤波器* 。 最常用的滤波器是 *线性* 滤波器。

1、`blur`均值滤波

文档地址：https://docs.opencv.org/4.x/d4/d86/group__imgproc__filter.html#ga8c45db9afe636703801b0b2e440fce37

```c++
blur	(	InputArray 	src,
            OutputArray  dst,
            Size 		ksize,
            Point 		anchor = Point(-1,-1),
            int 		borderType = BORDER_DEFAULT 
)	
```

`ksize`为滤波器大小。

`anchor`指定锚点位置(被平滑点)， 如果是负值，取核的中心为锚点。

`borderType`推断边缘像素，一般取默认值BORDER_DEFAULT。

2、`GaussianBlur`高斯滤波

```c++
GaussianBlur	(	InputArray 	src,
                    OutputArray 	dst,
                    Size 	ksize,
                    double 	sigmaX,
                    double 	sigmaY = 0,
                    int 	borderType = BORDER_DEFAULT 
)	
```

将输入数组的每一个像素点与 高斯内核*卷积，将卷积和当作输出像素值。

`sigmaX`x方向的标准方差。可设置为0让系统自动计算。

`sigmaY`y方向的标准方差。可设置为0让系统自动计算。

3、`medianBlur`中值滤波

```c++
medianBlur	(	InputArray 	src,
                OutputArray 	dst,
                int 	ksize 
)	
```

将图像的每个像素用邻域 (以当前像素为中心的正方形区域)像素的中值代替，中值值的是中间点的值，而不是平均值。

4、`bilateralFilter`双边滤波

```c++
bilateralFilter	(	InputArray 	src,
                    OutputArray 	dst,
                    int 	d,
                    double 	sigmaColor,
                    double 	sigmaSpace,
                    int 	borderType = BORDER_DEFAULT 
)	
```

执行双边滤波操作,类似于高斯滤波器，双边滤波器也给每一个邻域像素分配一个加权系数。 这些加权系数包含两个部分, 第一部分加权方式与高斯滤波一样，第二部分的权重则取决于该邻域像素与当前像素的灰度差值。

`d`像素的邻域直径。

`sigmaColor`颜色空间的标准方差。

`sigmaSpace`坐标空间的标准方差(像素单位)。

##### 图像翻转

```go
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
```

##### 图像阈值化/二值化

所谓阈值化，就是达到某个阈值的时候就被设置为另外一个值。

当灰度值大于一个值(阀值)时，让其成为一个大的值。 比如:灰度在0-255，当灰度小于128时赋值为0，大于128时赋值为255，即亮的地方更亮,暗的地方更暗。即实现了阀值分割，这样图像就黑白分明，对比度加大了。 阀值处理后使图象只有几种颜色如最通常的分为了黑白的二值图象。

```c++
threshold	(	InputArray 	src,
                OutputArray 	dst,
                double 	thresh,
                double 	maxval,
                int 	type 
)	
```

`thresh`比较直。

`maxval`在类型`THRESH_BINARY and THRESH_BINARY_INV`的时候会被用到。

`type`阈值化规则，可选值如下。

![image-20230627163547295](D:\dev\php\magook\trunk\server\md\img\image-20230627163547295.png)

`THRESH_BINARY`：二进制阈值，新的阈值产生规则可为：`value > thresh ? maxval : 0`

`THRESH_BINARY INV`：反二进制阈值，新的阈值产生规则可为：`value > thresh ? 0: maxval `

`THRESH_TRUNC`：截断阈值，新的阈值产生规则可为：`value > thresh ? thresh : value`

`THRESH_TOZERO`：阈值化为0，新的阈值产生规则可为：` value > thresh ? value : 0`

`THRESH_TOZERO_INV`：反阈值化为0，新的阈值产生规则可为：`value > thresh ? 0 : value`

```go
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
```

![image-20230627165935558](D:\dev\php\magook\trunk\server\md\img\image-20230627165935558.png)



##### 实现正交匹配追踪算法OMP

##### 计算协方差矩阵



https://www.cnblogs.com/denny402/category/716241.html