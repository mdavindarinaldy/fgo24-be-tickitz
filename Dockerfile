FROM golang:alpine AS build
WORKDIR /buildapp

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gotickitz main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=build /buildapp/gotickitz /app/

ENTRYPOINT [ "/app/gotickitz" ]