FROM golang:1.25.1-alpine3.22 AS build

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o votefin ./main.go



# FROM alpine:stable-slim
#
# COPY votefin /bin/votefin
#
# CMD ["/bin/votefin"]
