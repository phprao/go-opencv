QRCode二维码生成规范以及识别技术原理

1、摄像头

2、二维码扫描枪

3、光电扫描设备



##### QR码的特点和优势

以QR码（Quick Response Code）为例，QR码是由日本Dens0公司于1994年9月研制的一种矩阵式二维条码，它除了具有二维条码所具有的信息容量大、可靠性高、可表示汉字及图像多种信息、保密防伪性强等特点外，还具有能高速全方位识读、能有效表达汉字等主要特点。

每种码制有其特定的字符集，每个字符占有一定的宽度



二维码其实就是由很多0、1组成的数字矩阵。二维码是用某种特定的几何图形按一定规律在平面(二维方向上)分布的黑白相间的图形记录数据符号信息的，在代码编制上巧妙地利用构成“1”比特流的概念，使用若千个与二进制相对应的几何形体来计算机内部逻辑基础的0”表示文字数值信息，通过图象输入设备或光电扫描设备自动识读以实现信息自动处理: 它具有条码技术的一些共性:每种码制有其特定的字符集;每人字符占有一定的宽度，具有一定的校验功能等。同时还具有对不同行的信息自动识别功能、及处理图形旋转变化等特点。二维条码/二维码能够在横向和纵向两个方位同时表达信息，因此能在很小的面积内表达大量的信息



简单的说，二维码就是把你想表达的信息翻译成黑白两种小方块，然后填到这个大方块中。有点类似我们中学的答题卡，就是把我们的语言翻译成机器可识别的语言，说白了就是把数字.字母、汉字等信息通过特定的编码翻译成二进制0和1，一个0就是一个白色小方块，一个1就是一个黑色小方块。
当然这其中还有很多纠错码，假如需要编码的码字数据有100个，并且想对其中的一半，也就是50个码字进行纠错，则计算方法如下。纠错需要相当于码字2倍的符号，因此在这种情况下的数量为50个x2=100码字。因此，全部码字数量为200个，其中用作纠错的码字为50个，也就是说在这个二维码中，有25%的信息是用来纠错的，所以这也就解释了二维码即使缺了一点或者变皱了也一样能被识别。有些朋友可能会问，为什么每个二维码上都会有三个黑色大方块呢? 那就要涉及下面的内容:手机是如何识别二维码的。



二维码怎么被手机识别的?
由于不同颜色的物体，其反射的可见光的波长不同，白色物体能反射各种波长的可见光，黑色物体则吸收各种波长的可见光,所以当摄像头扫描黑白相间的二维码上时，手机利用点运算的茂直理论将采集到的图象变为二值图像,即对图像进行二值化处理,得到二值化图像后,对其进行膨胀运算,对膨胀后的图象进行边缘检测得到条码区域的轮廓

然后经过一项灰度值计算公式对图像进行二值化处理。得到一幅标准的二值化图像后，对该符号进行网格采样，对网格每一个交点上的图像像素取样，并根据闻值确定是深色“1”还是浅色“0”，从而得到二维码的原始二进制序列值，然后对这些数据进行纠错和译码，最后根据条码的逻辑编码规则把这些原始的数据转换成数据.
上文中我们提到的三个大黑方块起什么作用呢? 我们在使用手机扫描的时候无论是什么方向都能够正确识别二维码的内容，就是因为手机通过三个大黑方块识别出二维码正确的方向。上述识别过程说起来分很多步骤，但在实际生活中，也就是两三秒的工夫，非常快!

以下的讨论, 认定二维码是常见的 QR code.

摄像头首先获得图像,

对图像进行初步处理, 例如灰度, 锐化, 旋转等方式让图像的内容更加单纯;

然后根据 ISO/IEC 18004 的标准解码.

解码过程可以参考:



为程序员写的Reed-Solomon码

https://www.jianshu.com/p/8208aad537bb

https://zhuanlan.zhihu.com/p/542393490

https://blog.csdn.net/sinat_22510827/article/details/109648621

https://www.felix021.com/blog/read.php?entryid=2116&page=1&part=1



最小的二维码是21x21模块，最大的是177x177模块。尺寸称为Version。 21x21模块大小是Version 1，增加4，25x25模块大小是Version 2，依此类推。 177x177大小模块是Version 40。



纠错码的级别

最低的是L，它可以校准7%的字码。之后是可以校准15％的M，然后是可以校准25％的Q，最后是可以校准30％的H。一个二维码的容量取决于它的版本和错误纠正级别，以及编码的数据类型。二维码可以编码四种数据模式：数字，字符，字节或日文。（ps：实际上还有 Extended Channel Interpretation (ECI) mode、Structured Append mode 和 FNC1 mode，一般情况下用不到，所以作者也没有介绍。）



快速响应矩阵码QRCode国家标准 GB/T 18284-2000

全国标准信息公共服务平台：https://std.samr.gov.cn，看文档的地方都没有

中国标准在线服务网：https://www.spc.org.cn/，只能在线看

国家标准网：http://www.biaozhun8.cn/，在这个民间网址上找到了下载





##### QR码结构图

定位矫正图形

透视变换

识别和纠错

QR码图像预处理

图像灰度化

去噪

畸变矫正

二值化



微信小程序的圆形二维码，只能微信APP才能识别出来。
![在这里插入图片描述](https://img-blog.csdnimg.cn/fd7deac15fbc4784ac3d3a565f4c1b5c.png)




抖音的圆形二维码，只能抖音APP才能识别出来，抖音视频中禁止挂外链二维码，但是可以挂抖音码，

![在这里插入图片描述](https://img-blog.csdnimg.cn/ca8ea19914724322ba6776293c6650ce.png)




https://www.zhihu.com/search?type=content&q=%E4%BA%8C%E7%BB%B4%E7%A0%81%E7%9A%84%E5%8E%9F%E7%90%86

https://zhuanlan.zhihu.com/p/25423714

http://www.360doc.com/showweb/0/0/1067328716.aspx

https://blog.csdn.net/search_129_hr/article/details/120796256

https://coolshell.cn/articles/10590.html

https://blog.csdn.net/bemy1008/article/details/82886915

https://blog.csdn.net/weiwei9363/article/details/81180966

[《二维码的生成细节和原理》](http://blog.csdn.net/hk_5788/article/details/50839790)
[《QR Code Tutorial》](https://www.thonky.com/qr-code-tutorial/)
[《Hello World!》—— 知乎专栏文章](https://zhuanlan.zhihu.com/p/21463650)
[《为程序员写的Reed-Solomon码解释》](http://www.jianshu.com/p/8208aad537bb)

https://www.jianshu.com/p/3cf1862552f8
https://blog.csdn.net/u012611878/article/details/53167009
https://blog.51cto.com/jsxyhelu2017/5972864
https://www.open-open.com/lib/view/open1383022297577.html


golang二维码操作 github.com/skip2/go-qrcode
```go
package main

import (
	"image"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/skip2/go-qrcode"
)

func main() {
	qrcodeImage()
}

func qrcodeImage() {
	qr, err := qrcode.New("中华人民共和国", qrcode.Medium)
	if err != nil {
		panic(err)
	}
	showWindow(qr.Image(256))
}

func qrcodePng() {
	var pngBytes []byte
	pngBytes, _ = qrcode.Encode("https://example.org", qrcode.Medium, 256)

	f, _ := os.Create("code1.png")
	defer f.Close()

	f.Write(pngBytes)
}

func showWindow(image image.Image) {
	myApp := app.New()
	myWindow := myApp.NewWindow("二维码")

	img2 := canvas.NewImageFromImage(image)
	img2.FillMode = canvas.ImageFillOriginal

	myWindow.SetContent(img2)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

```

golang调用摄像头opencv


golang中用来识别和解析二维码的包：github.com/makiuchi-d/gozxing，此包实现了 zxing 的算法，效果较好，在 opencv 解析不出的时候它能解析出来，市面上大都是基于 zxing 来实现的。