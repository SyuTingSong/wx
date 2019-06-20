FROM golang:1-alpine AS build
COPY . /build
WORKDIR /build
RUN go mod init && go build -o wx .

FROM alpine:3.9
COPY --from=build /build/wx /app/wx
ENV GIN_MODE=release
ENTRYPOINT ["/app/wx"]
