FROM golang:bullseye AS builder
ADD ffmpeg.go /build/ffmpeg.go
RUN \
cd /build && \
go build -ldflags="-s -w" ffmpeg.go

FROM emby/embyserver_arm64v8:4.9.1.12
COPY --from=builder /build/ffmpeg /bin/my-ffmpeg
RUN \
mv /bin/ffmpeg /bin/_ffmpeg && \
mv /bin/my-ffmpeg /bin/ffmpeg
