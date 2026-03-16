FROM golang:1.25.1-alpine AS build

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o votefin .

FROM golang:1.25.1-alpine AS db

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

RUN go install github.com/pressly/goose/v3/cmd/goose@latest   

COPY --from=build /app/sql/schema/ ./

RUN goose sqlite3 ./votefin.db up

FROM alpine:3.23.3 AS prod

ENV DEV_MODE="PROD"
WORKDIR /app

COPY --from=build /app/votefin .
COPY --from=build /app/assets/styles.css ./assets/styles.css
COPY --from=db /app/votefin.db ./data/votefin.db
COPY --from=build /app/data/ ./data/


EXPOSE 8080
CMD ["./votefin"]

