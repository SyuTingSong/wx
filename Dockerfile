FROM golang:1-alpine AS build
COPY . /build
WORKDIR /build
RUN apk --no-cache add git && \
    go mod download && \
    go build -o wx .

FROM alpine:3.9
COPY --from=build /build/wx /app/wx
ENV GIN_MODE=release
ENTRYPOINT ["/app/wx"]

