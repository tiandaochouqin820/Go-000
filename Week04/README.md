1、安装pb

	go get -u github.com/golang/protobuf/protoc-gen-go

2、安装go-zero和goctl工具

	go get -u github.com/tal-tech/go-zero
	go get -u github.com/tal-tech/go-zero/tools/goctl

3、创建工程目录shorturl和api目录shorturl/api

    mkdir -p shorturl/api

4、在shorturl目录中执行

	go mod init shorturl

	添加依赖项

	module shorturl

	go 1.15
	
	require (
	  github.com/golang/mock v1.4.3
	  github.com/golang/protobuf v1.4.2
	  github.com/tal-tech/go-zero v1.1.1
	  golang.org/x/net v0.0.0-20200707034311-ab3426394381
	  google.golang.org/grpc v1.29.1
	)

5、在shorturl/api目录下利用goctl生成api

	goctl api -o shorturl.api

6、利用goctl生成API代码

	goctl api go -api shorturl.api -dir .

7、切换至shorturl目录，创建pb模板

	mkdir -p rpc/transform
	cd rpc/transform/
	goctl rpc template -o transform.proto

8、生成rpc代码

	goctl rpc proto -src transform.proto -dir .

9、切换至shorturl目录，创建model目录，创建shorturl.sql文件

	cd ../../
	mkdir -p rpc/transform/model

10、生成DB所需代码

	goctl model mysql ddl -c -src shorturl.sql -dir .