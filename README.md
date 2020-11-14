# Hello Fudan Go!

一个Go语言版本的平安复旦，支持多用户同时签到

## 使用方法

* clone 代码并运行二进制文件

  ```bash
  git clone https://github.com/kagaya85/HelloFudanGo.git
  cd HelloFudanGO/build
  ./helloFudanGo
  ```

* 账户信息保存在项目目录下的`accounts.txt`，如需要同时多用户签到，只需要编辑该文件，按相同格式添加账户密码（空格分隔，每行一个用户）

* 编译源码

  ```bash
  # go 1.14
  go install
  go build -o build/helloFudanGo
  ```


* **If you find any mistakes in code, issues or pull request is always WELCOME!**

