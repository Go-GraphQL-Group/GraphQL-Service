#源镜像
FROM golang:latest
#作者
MAINTAINER Razil "Go-GraphQL-Group"
#设置工作目录
WORKDIR $GOPATH/src/github.com/Go-GraphQL-Group/GraphQL-Service
#将服务器的go工程代码加入到docker容器中
ADD . $GOPATH/src/github.com/Go-GraphQL-Group/GraphQL-Service
#go构建可执行文件
RUN go build .
#暴露端口
EXPOSE 9090
#最终运行docker的命令
ENTRYPOINT  ["./main.go"]