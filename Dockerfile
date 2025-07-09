FROM golang:alpine AS build
WORKDIR /buildapp

COPY . .

RUN go build -o gotickitz main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=build /buildapp/gotickitz /app/gotickitz

ENTRYPOINT [ "/app/gotickitz" ]