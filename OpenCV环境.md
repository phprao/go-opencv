##### 演化流程

`OpenCV  --> 传统的机器学习SIFT --> 神经网络CNN`

计算机视觉也由诸如梯度方向直方图（Histogram of Gradient， HOG）以及尺度不变特征变换（Scale-Invariant Feature Transform， SIFT）等传统的手办特征（Hand-Crafted Feature）与浅层模型的组合逐渐转向了以卷积神经网络（Convolutional Neural Network， CNN）为代表的深度学习模型。

OpenCV的全称是Open Source Computer Vision Library，是一个跨平台的计算机视觉库。OpenCV是由英特尔公司发起并参与开发，以BSD许可证授权发行，可以在商业和研究领域中免费使用。OpenCV可用于开发实时的图像处理、计算机视觉以及模式识别程序。该程序库也可以使用英特尔公司的IPP进行加速处理。


#### golang配置OpenCV

##### 1、关于OpenCV

https://github.com/opencv/opencv/tree/4.7.0

https://opencv.org/releases/

https://docs.opencv.org/4.7.0/  当前稳定版本 4.7.0

opencv论坛 https://forum.opencv.org/


##### 2、安装GoCV

官网 https://gocv.io/

对应的仓库 `https://github.com/hybridgroup/gocv`  基于`OpenCV-4.7.0`，大量示例在`cmd`目录下。

先安装`MinGW-W64 v8.1.0`和`CMake`并添加到环境变量。

下载包 `go get -u -d gocv.io/x/gocv`

进入到`GOPATH`目录下找到`gocv`，使用PowerShell（不能是普通的cmd窗口）执行`win_build_opencv.cmd`来下载安装`opencv`，默认安装到了`C:\opencv`。

需要一些时间。

然后将`C:\opencv\build\install\x64\mingw\bin`添加到环境变量。

验证安装是否成功，进入到`GOPATH`目录下找到`gocv`，执行`go run cmd\version\main.go`。输出如下说明安装成功

```bash
gocv version: 0.33.0
opencv lib version: 4.7.0
```

注意，opencv的功能是大量依赖于你本地环境的，在安装opencv的时候会开到大量的`not found`，然后它就将对应的功能设置为`off`，于是你在使用的时候就会报错，如果的确需要某功能，需要先装上对应的依赖，然后重新编译即可。

在安装的时候，有的系统会报错`D:\Program`不是一个可运行程序，那是因为MinGW被安装在了`D:\Program Files`目录下，在命令行模式下，这个空格被解析为前面是程序，后面是参数，因此报错，我的解决方式是修改`win_build_opencv.cmd`文件，拉到后面，将`mingw32-make`写成了绝对路径，并加上了双引号，如下：`D:\"Program Files"\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64\bin\mingw32-make`。或者重装`MinGW`。


学习文档：
https://blog.csdn.net/qq_15698613/category_9292368.html
https://github.com/hybridgroup/gocv/cmd


#### OpenCV概述

各个模块的功能：

`calib3d`
其实就是就是Calibration（校准）加3D这两个词的组合缩写。这个模块主要是相机校准和三维重建相关的内容。基本的多视角几何算法，单个立体摄像头标定，物体姿态估计，立体相似性算法，3D信息的重建等等。

`core`
核心功能模块，包含如下内容：

- OpenCV基本数据结构
- 动态数据结构
- 绘图函数
- 数组操作相关函数
- 辅助功能与系统函数和宏
- 与OpenGL的互操作

`imgproc`

Image和Processing这两个单词的缩写组合。图像处理模块，这个模块包含了如下内容：

- 线性和非线性的图像滤波

- 图像的几何变换
- 其它（Miscellaneous）图像转换
- 直方图相关
- 结构分析和形状描述
- 运动分析和对象跟踪
- 特征检测
- 目标检测等内容

`features2d`

也就是Features2D， 2D功能框架 ，包含如下内容：

- 特征检测和描述

- 特征检测器（Feature Detectors）通用接口
- 描述符提取器（Descriptor Extractors）通用接口
- 描述符匹配器（Descriptor Matchers）通用接口
- 通用描述符（Generic Descriptor）匹配器通用接口
- 关键点绘制函数和匹配功能绘制函数

`flann`

 Fast Library for Approximate Nearest Neighbors，高维的近似近邻快速搜索算法库，包含两个部分

- 快速近似最近邻搜索
- 聚类

`highgui`

也就是high gui，高层GUI图形用户界面，包含媒体的I / O输入输出，视频捕捉、图像和视频的编码解码、图形交互界面的接口等内容

`ml`


Machine Learning，机器学习模块， 基本上是统计模型和分类算法，包含如下内容：

- 统计模型 （Statistical Models）
- 一般贝叶斯分类器 （Normal Bayes Classifier）
- K-近邻 （K-NearestNeighbors）
- 支持向量机 （Support Vector Machines）
- 决策树 （Decision Trees）
- 提升（Boosting）
- 梯度提高树（Gradient Boosted Trees）
- 随机树 （Random Trees）
- 超随机树 （Extremely randomized trees）
- 期望最大化 （Expectation Maximization）
- 神经网络 （Neural Networks）
- MLData

`objdetect`

目标检测模块，包含Cascade Classification（级联分类）和Latent SVM这两个部分。

`photo`

也就是Computational Photography，包含图像修复和图像去噪两部分

`stitching`

images stitching，图像拼接模块，包含如下部分：

- 拼接流水线
- 特点寻找和匹配图像
- 估计旋转
- 自动校准
- 图片歪斜
- 接缝估测
- 曝光补偿
- 图片混合

`video`


视频分析组件，该模块包括运动估计，背景分离，对象跟踪等视频处理相关内容。





