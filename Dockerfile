FROM golang as build

#COPY web/public  /usr/local/go/src/lottery/web/public
COPY out/remind /usr/local/go/src/remind/out/remind

WORKDIR /usr/local/go/src/remind/
ADD  out/remind remind
# 使用C语言版本的GO编译器，参数配置为0的时候就关闭C语言版本的编译器了 生成64位linux的可执行程序
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server
COPY --from=build /usr/local/go/src/remind/remind /usr/bin/remind
RUN chmod +x /usr/bin/remind

ENTRYPOINT ["remind"]
