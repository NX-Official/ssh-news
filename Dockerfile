FROM golang:latest AS go-build-stage

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM ubuntu:latest
RUN apt-get update && apt-get install -y bash bash-completion openssl ca-certificates
RUN update-ca-certificates
ENV TERM=xterm-256color

WORKDIR /app

COPY --from=go-build-stage /go/src/app/main .

CMD ["./main"]
