OCR（Optical character recognition，光学字符识别）是一种将图像中的手写字或者印刷文本转换为机器编码文本的技术。

[tesseract-ocr](https://github.com/tesseract-ocr/tesseract) 是由Google开发，支持[100多种语言](https://tesseract-ocr.github.io/tessdoc/Data-Files-in-different-versions.html)

文档 tessdoc：

[https://tesseract-ocr.github.io/](https://tesseract-ocr.github.io/)

[https://tesseract-ocr.github.io/tessdoc/Installation.html](https://tesseract-ocr.github.io/tessdoc/Installation.html)

[https://github.com/tesseract-ocr/tessdoc](https://github.com/tesseract-ocr/tessdoc)

Windows Tesseract下载地址：[https://digi.bib.uni-mannheim.de/tesseract/](https://digi.bib.uni-mannheim.de/tesseract/)

选择 [tesseract-ocr-w64-setup-v5.0.1.20220118.exe](https://digi.bib.uni-mannheim.de/tesseract/tesseract-ocr-w64-setup-v5.0.1.20220118.exe)

![image-20230713160559046](D:\dev\php\magook\trunk\server\md\img\image-20230713160559046.png)

勾选上`Additional...`会下载训练数据，安装到 `D:\Tesseract-OCR`，将`D:\Tesseract-OCR`添加到环境变量。

```bash
C:\Users\Administrator.DESKTOP-TPJL4TC>tesseract
Usage:
  tesseract --help | --help-extra | --version
  tesseract --list-langs
  tesseract imagename outputbase [options...] [configfile...]

OCR options:
  -l LANG[+LANG]        Specify language(s) used for OCR.
NOTE: These options must occur before any configfile.

Single options:
  --help                Show this help message.
  --help-extra          Show extra help for advanced users.
  --version             Show version information.
  --list-langs          List available languages for tesseract engine.
  
C:\Users\Administrator.DESKTOP-TPJL4TC>tesseract --version
tesseract v5.0.1.20220118
 leptonica-1.78.0
  libgif 5.1.4 : libjpeg 8d (libjpeg-turbo 1.5.3) : libpng 1.6.34 : libtiff 4.0.9 : zlib 1.2.11 : libwebp 0.6.1 : libopenjp2 2.3.0
 Found AVX2
 Found AVX
 Found FMA
 Found SSE4.1
 Found libarchive 3.5.0 zlib/1.2.11 liblzma/5.2.3 bz2lib/1.0.6 liblz4/1.7.5 libzstd/1.4.5
 Found libcurl/7.77.0-DEV Schannel zlib/1.2.11 zstd/1.4.5 libidn2/2.0.4 nghttp2/1.31.0
```

查看支持的语言包

```bash
tesseract --list-langs
```

如果忘记勾选了训练数据，也可以单独下载 [https://digi.bib.uni-mannheim.de/tesseract/tessdata_fast/](https://digi.bib.uni-mannheim.de/tesseract/tessdata_fast/)，放在`D:\Tesseract-OCR\tessdata`目录下

识别图片中的文字，默认只能识别英文和数字

```bash
tesseract 图片地址 存放识别结果的文本文件路径
比如
tesseract D:\dev\php\magook\trunk\server\go-opencv\detectcarband\licence_plate.jpg D:\dev\php\magook\trunk\server\go-opencv\detectcarband\licence_plate
```

如果要识别中文，那就需要加上语言包名称

```bash
tesseract D:\dev\php\magook\trunk\server\go-opencv\detectcarband\licence_plate.jpg D:\dev\php\magook\trunk\server\go-opencv\detectcarband\licence_plate -l chi_sim
```

其实并不算很准，比如如下车牌

![image-20230713165033900](./imgs/image-20230713165033900.png)

识别结果是`外.730V7`

```bash
C:\Users\Administrator.DESKTOP-TPJL4TC> tesseract --help-extra
Usage:
  tesseract --help | --help-extra | --help-psm | --help-oem | --version
  tesseract --list-langs [--tessdata-dir PATH]
  tesseract --print-fonts-table [options...] [configfile...]
  tesseract --print-parameters [options...] [configfile...]
  tesseract imagename|imagelist|stdin outputbase|stdout [options...] [configfile...]

OCR options:
  --tessdata-dir PATH   Specify the location of tessdata path.
  --user-words PATH     Specify the location of user words file.
  --user-patterns PATH  Specify the location of user patterns file.
  --dpi VALUE           Specify DPI for input image.
  --loglevel LEVEL      Specify logging level. LEVEL can be
                        ALL, TRACE, DEBUG, INFO, WARN, ERROR, FATAL or OFF.
  -l LANG[+LANG]        Specify language(s) used for OCR.
  -c VAR=VALUE          Set value for config variables.
                        Multiple -c arguments are allowed.
  --psm NUM             Specify page segmentation mode.
  --oem NUM             Specify OCR Engine mode.
NOTE: These options must occur before any configfile.

Page segmentation modes:
  0    Orientation and script detection (OSD) only.
  1    Automatic page segmentation with OSD.
  2    Automatic page segmentation, but no OSD, or OCR. (not implemented)
  3    Fully automatic page segmentation, but no OSD. (Default)
  4    Assume a single column of text of variable sizes.
  5    Assume a single uniform block of vertically aligned text.
  6    Assume a single uniform block of text.
  7    Treat the image as a single text line.
  8    Treat the image as a single word.
  9    Treat the image as a single word in a circle.
 10    Treat the image as a single character.
 11    Sparse text. Find as much text as possible in no particular order.
 12    Sparse text with OSD.
 13    Raw line. Treat the image as a single text line,
       bypassing hacks that are Tesseract-specific.

OCR Engine modes:
  0    Legacy engine only.
  1    Neural nets LSTM engine only.
  2    Legacy + LSTM engines.
  3    Default, based on what is available.

Single options:
  -h, --help            Show minimal help message.
  --help-extra          Show extra help for advanced users.
  --help-psm            Show page segmentation modes.
  --help-oem            Show OCR Engine modes.
  -v, --version         Show version information.
  --list-langs          List available languages for tesseract engine.
  --print-fonts-table   Print tesseract fonts table.
  --print-parameters    Print tesseract parameters.
```

golang使用tesseract-ocr：[gosseract](https://github.com/otiai10/gosseract)

sesseract-ocr 使用教程：https://blog.csdn.net/u010698107/article/details/121736386

sesseract-ocr 识别不准怎么办