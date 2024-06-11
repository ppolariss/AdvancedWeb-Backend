FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache gcc\
        tzdata \
        gcc \
        g++ &&  \
    go mod download

COPY . .

RUN go build -o AdvancedWebpj

FROM alpine

WORKDIR /app

COPY --from=builder /app/AdvancedWebpj .
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# ENV DB_URL root:root@tcp(host.docker.internal:3306)/advanced_web?charset=utf8mb4&parseTime=True&loc=Local
ENV TZ=Asia/Shanghai

EXPOSE 8080

ENTRYPOINT ["./AdvancedWebpj"]