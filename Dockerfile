FROM golang:latest as builder

COPY . /app
WORKDIR /app/cmd


RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod download
RUN go build -ldflags "-w -s" -gcflags "-N -l" -o main .

# RUN
FROM alpine:latest
COPY --from=builder /app/cmd/main /app/main
ENTRYPOINT ["/app/main"]