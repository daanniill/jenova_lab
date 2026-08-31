FROM golang:1.27-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY service ./service

RUN go build -o /service ./service

FROM alpine:latest

COPY --from=build /service /service

EXPOSE 8080

CMD ["/service"]