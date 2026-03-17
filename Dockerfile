FROM golang:1.25.1-alpine AS build

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o votefin .

FROM alpine:3.23.3 AS prod

ENV DEV_MODE="PROD"
WORKDIR /app

COPY --from=build /app/votefin .
COPY --from=build /app/assets/styles.css ./assets/styles.css
COPY --from=build /app/data/ ./data/


EXPOSE 8080
CMD ["./votefin"]

