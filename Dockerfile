# dev
FROM golang:1.18-bullseye AS dev
WORKDIR /work/yatter-backend-go
RUN go install github.com/cosmtrek/air@v1.42.0

COPY ./ ./
RUN make mod build-linux

# release
FROM alpine AS release
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
COPY --from=dev /work/yatter-backend-go/build/yatter-backend-go-linux-amd64 /usr/local/bin/yatter-backend-go
EXPOSE 8080
ENTRYPOINT ["yatter-backend-go"]
