# FileHive  
像快递柜一样存取文件

## 项目说明  
FileHive的设计灵感来自于Filecodebox。通过取件码和取件密码，用户可以方便地获取存储在系统中的文件。该项目基于`fiber`框架和`BadgerDB`数据库开发。

## 功能展示  

![GitHub图像](/1.png)  
![GitHub图像](/2.png)  
![GitHub图像](/3.png)  
![GitHub图像](/4.png)  
![GitHub图像](/5.png)  

## 构建方法  

1. 克隆项目代码到本地：
    ```bash
    git clone https://github.com/your-username/FileHive.git
    cd FileHive
    ```

2. 安装Go依赖：
    ```bash
    go mod tidy
    ```

3. 编译项目：
    ```bash
    go build .
    ```

## 运行方法  

1. 在Windows或Linux系统中，直接运行编译后的二进制文件：
    ```bash
    ./filehive
    ```

2. 打开浏览器，访问[http://127.0.0.1:3000](http://127.0.0.1:3000)，即可使用FileHive系统。

