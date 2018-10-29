# 变量
# ARG MYGOPATH="/Users/elvin/Workspace/Golang"

# 基于golang image
FROM golang:1.9

# 作者信息
MAINTAINER gopher team

# 将当前 目录文件 拷贝到 container 中的目录
COPY . /go/src/igen/msdemo

# 将当前 目录文件 拷贝到 container 中的目录
# ADD vs COPY:
#   基本功能一样
#   ADD 可以添加一个URL
#   ADD 可以添加tar等，并会自动解压, 并删除原tar文件
#ADD ${MYGOPATH}/bin/logf /go/bin/

# contanier 需要对外开放的端口
EXPOSE 8081

# 环境变量
# ENV <key> <value>

# 共享目录
# VOLUME ["samepath"]

# WORKDIR
# WORKDIR指令用于设置Dockerfile中的RUN、CMD和ENTRYPOINT指令执行命令的工作目录(默认为/目录)。
# 可以出现多次，如果使用相对路径则为相对于WORKDIR上一次的值，例如WORKDIR /data，WORKDIR logs，RUN pwd最终输出的当前目录是/data/logs
WORKDIR /go/src/igen/msdemo/vendor/igen/lib/logger/logf
RUN go install

WORKDIR /go/src/igen/msdemo/vendor/github.com/silenceper/gowatch
RUN go install

# RUN
# contanier 构建时需要执行的命令
#RUN go install

# contanier 运行时要执行的命令
WORKDIR /go/src/igen/msdemo
#CMD ["go", "run", "main.go", "|", "logf"]
CMD ["gowatch"]