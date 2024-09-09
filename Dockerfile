FROM golang:1.22 as build

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY main.go .
COPY mime_types.json .

RUN go build -o /app/main /app/main.go

FROM ubuntu:24.04 as prod

WORKDIR /app

COPY --from=build /app/main /app/main
COPY mime_types.json .

CMD ["/app/main"]
