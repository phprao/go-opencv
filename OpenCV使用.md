`func NewWindow(name string) *Window`

创建窗口，可以创建多个，但是名称要不一样。

opencv中的GUI是由`HighGUI`这个组件提供的。

`func (w *Window) WaitKey(delay int) int`

用于阻塞当前线程，否则窗口会一闪而过。如果参数为0就会一直阻塞，直到有按键事件；如果大于0，就会等待相应的毫秒然后解除阻塞，或者有按键事件就会立即解除阻塞。返回值为按键的code，如果没有按键则返回值为-1。如果有多个窗口，只需要调用一次 WaitKey 会同时生效。

`func IMRead(name string, flags IMReadFlag) Mat`

加载图片文件。常用的读取模式：1、`IMReadColor`读取出来的为彩色图片，拥有三个通道，opencv中的通道顺序为`BGR`，而不是`RGB`；2、`IMReadGrayScale`读取出来的为黑白照片，只有一个通道。支持的图片格式很多，比如常用的`jpg,jpeg,png,bmp,webp,awf`，具体可以看文档，不支持`gif`。

`gocv.Mat`矩阵结构，它由两部分组成，一个是 matrix header，这部分大小是固定的，包含矩阵的大小，存储的方式，矩阵存储的地址等等，另一部分是一个指针，指向了矩阵元素。这类似于切片结构，因此应该使用`Clone()`或者`CopyTo()`函数来实现赋值和拷贝，否则对 Mat 的操作会相互影响，因为普通的赋值只复制了header，而底层的数据指针是一样的，这一做的好处是避免了矩阵的重复分配内存。

```c++
// https://github.com/opencv/opencv/blob/4.7.0/modules/core/include/opencv2/core/mat.hpp

class CV_EXPORTS Mat
{
	...
 
    int dims;  /*数据的维数*/
    int rows,cols; /*行和列的数量;数组超过2维时为(-1，-1)*/
    uchar *data;   /*指向数据*/
    int * refcount;   /*指针的引用计数器; 阵列指向用户分配的数据时，指针为 NULL*/

    int flags; // 重点讲解
    //! the matrix dimensionality, >= 2，数据的维数
    int dims;
    //! the number of rows and columns or (-1, -1) when the matrix has more than 2 dimensions
    int rows, cols;
    //! pointer to the data，指向数据
    uchar* data;

    //! helper fields used in locateROI and adjustROI
    const uchar* datastart;
    const uchar* dataend;
    const uchar* datalimit;

    //! custom allocator
    MatAllocator* allocator;
    //! and the standard allocator
    static MatAllocator* getStdAllocator();
    static MatAllocator* getDefaultAllocator();
    static void setDefaultAllocator(MatAllocator* allocator);

    //! internal use method: updates the continuity flag
    void updateContinuityFlag();

    //! interaction with UMat
    UMatData* u;

    MatSize size;
    MatStep step;
 
	...
 
};
```

上面的`int flags`属性，占用32位，从低位到高位

```bash
0-2位代表depth即数据类型（如CV_8U），OpenCV的数据类型共7类，故只需3位即可全部表示。

3-11位代表通道数channels，因为OpenCV默认最大通道数为512，故只需要9位即可全部表示，可参照下面求通道数的部分。

0-11位共同代表type即通道数和数据类型（如CV_8UC3）

12-13位暂没发现用处，也许是留着后用，待发现了再补上。

14位代表Mat的内存是否连续，一般由creat创建的mat均是连续的，如果是连续，将加快对数据的访问。

15位代表该Mat是否为某一个Mat的submatrix，一般通过ROI以及row()、col()、rowRange()、colRange()等得到的mat均为submatrix。

16-31代表magic signature，暂理解为用来区分Mat的类型，如Mat和SparseMat
```

`MatType`矩阵的类型，包含了元素的数据类型depth以及通道个数，比如`MatTypeCV8UC3`代表`8 uint and channel 3`。depth的取值是从0到6。

`ElemSize()`返回矩阵中一个元素占用的字节数，比如`MatTypeCV8UC3`每个元素占用字节数为3。

```go
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
```

`Mat.reshape()`只是改变了矩阵的形状，改变前后切片内的值的个数不变，即`row1 * col1 * channel1 == row2 * col2 * channel2`；而`gocv.resize()`是改变图像的宽高，并不保证切片内的值的个数不变。

```go
func Resize(src Mat, dst *Mat, sz image.Point, fx, fy float64, interp InterpolationFlags)
```

`zs`：目标图像的size，gocv中通过`image.Pioint{}`来传递宽高，不要奇怪。

`fx, fy`：分别是X方向和Y方向的缩放比例，如果指定了`fx,fy`那么就不用传`sz`。同样的，如果指定了`sz`那么久不用传`fx,fy`。

`interp`：指定使用那种插值算法，因为图片的缩放效果跟使用那种插值算法有关，默认的是线性插值 Linear，即 `gocv.InterpolationLinear`。

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



##### 图像腐蚀与膨胀

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

最典型的一个应用场景就是在你`二值化`后，你的目标和背景抠的不是很干净，比如里面还有一些干扰点或者干扰线，可以试试两个操作，有时候效果出奇的好。

一般是对二值化后的图像来进行腐蚀操作，这样就是图形更加分明。处理图形周围的毛刺。每执行一次腐蚀操作，都会把图像多腐蚀一圈。

核的大小也会影响腐蚀的效果，核越小腐蚀的越多。



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

在opencv中，可能使用blur对图像进行平滑处理，这种方法就是最简单的求平均数。*平滑* 也称 *模糊*, 是一项简单且使用频率很高的图像处理方法。平滑处理的用途有很多， 但是在很多地方我们仅仅关注它减少噪声的功用。平滑处理时需要用到一个 *滤波器* 。 最常用的滤波器是 *线性* 滤波器。

滤波的原理和腐蚀膨胀类似，都是由“核”来遍历计算。核的边长一般设置为奇数。

1、`blur`均值滤波

核所覆盖的像素点的均值作为核中心处的颜色值。

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

将输入数组的每一个像素点与 高斯内核*卷积，将卷积和当作输出像素值。因为核覆盖的面积上，各个像素点距离核中心的距离不一样，导致颜色的贡献程度不一样，高斯滤波就是给不同的点不同的权重最终计算出核中心点的颜色值。

`sigmaX`x方向的标准方差。可设置为0让系统自动计算。

`sigmaY`y方向的标准方差。可设置为0让系统自动计算。

3、`medianBlur`中值滤波

```c++
medianBlur	(	InputArray 	src,
                OutputArray 	dst,
                int 	ksize 
)	
```

将图像的每个像素用邻域 (以当前像素为中心的正方形区域)像素的中值代替，中值值的是从大到小排序后中间的值，而不是平均值。

中值滤波非常适合处理细小的噪音点颗粒。

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

**关于一屏多图**

上面的示例中我们经常需要开启多个窗口来展示多个图片，比较麻烦，其实也可以在一个窗口中同时展示多个图片，其原理就是将多个图片编排到一个大的`Mat`中去，需要用到`Mat.Regin()`，得到的新Mat和原来的Mat是引用关系，我们将图片`CopyTo`到新Mat里去即可。需要注意的是，此方法要求这些图片的通道数一致，尺寸的话不做要求，可以`resize`之后再来编排。

```go
// imgCols 一行展示几个图片
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
```

![image-20230629151923803](D:\dev\php\magook\trunk\server\md\img\image-20230629151923803.png)

##### 读取GIF图像

opencv中无法读取gif图像，这是由于license原因。转而使用 videocapture 或者第三方的 PIL 库（Python），但是其实Golang的基础库`image`中就有读取gif图像的。于是一个简单的示例如下

```go
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
```

这里只会播放一遍gif图像，我们还可以解析gif中的LoopCount来增加循环播放的逻辑。

##### 读取mp4视频文件

首先要确保cmake安装的时候成功安装了`opencv_ffmpeg_64.dll and opencv_ffmpeg.dll`依赖，否则在调用`gocv.VideoCaptureFile`或者`gocv.OpenVideoCapture`的时候会报错`Error opening file: showimage/video1.mp4`。

打开opencv编译安装的路径下`C:\opencv\build\lib`，的确没找到这两依赖，那怎么办呢？

opencv在编译的时候会首先查找当前系统有没有安装ffmpeg，如果没有安装才会去下载安装，但是可能是在下载的时候失败了，所以就没有安装这个依赖，下载失败的日志可以在`opencv/build/CMakeDownloadLog.txt`找到，因此，我打开了梯子软件来重新编译opencv就可以了，不用梯子的话我还没试过怎么解决。

读取视频文件使用`gocv.VideoCaptureFile(filename)`或者`gocv.OpenVideoCapture(filename)`，然后逐帧处理

```go
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
```

其实也可以使用`ReadAndShowVideo`函数来读取GIF图像，但是不如`ReadAndShowGIF`控制的更细致。



https://github.com/BtbN/FFmpeg-Builds/releases

https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip

https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl-shared.zip



##### canny的原理





##### [opencv实现正交匹配追踪算法OMP](https://www.cnblogs.com/denny402/p/4676729.html)

##### [opencv2学习：计算协方差矩阵](https://www.cnblogs.com/denny402/p/5011456.html)

[opencv3中的机器学习算法之：EM算法](https://www.cnblogs.com/denny402/p/5036288.html)

[在opencv3中实现机器学习算法之：利用最近邻算法（knn)实现手写数字分类](https://www.cnblogs.com/denny402/p/5033898.html)

[在opencv3中的机器学习算法练习：对OCR进行分类](https://www.cnblogs.com/denny402/p/5032839.html)

[在opencv3中实现机器学习之：利用逻辑斯谛回归（logistic regression)分类](https://www.cnblogs.com/denny402/p/5032490.html)

[在opencv3中的机器学习算法](https://www.cnblogs.com/denny402/p/5032232.html)

[在opencv3中实现机器学习之：利用正态贝叶斯分类](https://www.cnblogs.com/denny402/p/5031613.html)

[在opencv3中进行图片人脸检测](https://www.cnblogs.com/denny402/p/5031181.html)

[在opencv3中利用SVM进行图像目标检测和分类](https://www.cnblogs.com/denny402/p/5020551.html)

[在opencv3中实现机器学习之：利用svm(支持向量机)分类](https://www.cnblogs.com/denny402/p/5019233.html)

[在matlab和opencv中分别实现稀疏表示](https://www.cnblogs.com/denny402/p/5016530.html)



https://www.cnblogs.com/denny402/category/716241.html

https://blog.csdn.net/youcans/article/details/125112487