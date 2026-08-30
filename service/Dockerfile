FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY service ./service

RUN go build -o /service ./service

FROM alpine:latest

COPY --from=build /service /service

EXPOSE 8080

CMD ["/service"]