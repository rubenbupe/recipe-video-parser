FROM golang:alpine AS build

# Install git and make
RUN apk add --update git make
WORKDIR /go/src/github.com/rubenbupe/recipe-video-parser
COPY . .
RUN make build

# Final image with shell and curl
FROM alpine:latest
WORKDIR /go
# Install curl, python3, pip, ffmpeg, gallery-dl, and yt-dlp
RUN apk add --no-cache curl python3 py3-pip ffmpeg \
    && pip install --no-cache-dir gallery-dl yt-dlp
COPY --from=build /go/src/github.com/rubenbupe/recipe-video-parser/bin/api /go/bin/api
COPY --from=build /go/src/github.com/rubenbupe/recipe-video-parser/bin/cli /go/bin/cli

ENTRYPOINT ["/go/bin/api"]