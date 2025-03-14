FROM golang:1.24-alpine
ENV GOPROXY=https://goproxy.cn

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]
