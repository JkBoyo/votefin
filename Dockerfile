FROM golang:1.25.1 AS build

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o votefin .

FROM alpine:3.23.3 AS prod

WORKDIR /app

COPY --from=build /app/votefin .
COPY --from=build /app/assets/styles.css ./assets/styles.css
COPY --from=build /app/data/ ./assets/data/

EXPOSE 8080
# CMD ["./votefin"]

