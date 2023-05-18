## 准备工作
```bash
brew install protobuf
protoc --version #最新版本
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH=$PATH:$GOPATH/bin
```
## 下载demo
```bash
git clone https://github.com/grpc-example/resolver.git
cd resolver
go mod tidy
```
## 启动两个服务端
```bash
go run server.go --port=50001
go run server.go --port=50002
```
## 启动客户端
```bash
go run client.go
curl http://localhost:8081/hello
```
多次curl调用，可以看到client的请求轮询发送到服务端
## 说明
client上启动了一个http服务，做为流量入口，然后通过grpcClient去请求grpc服务。
- 这里使用了grpc的manual(手动模式)包，自定义ds的scheme
- m.InitialState 时定义了两个服务的地址
- m.UpdateState 10s后更新成不存在的地址，所以后面请求就失败了