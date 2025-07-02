FROM golang:alpine AS build

# Install git and make
RUN apk add --update git make
WORKDIR /go/src/github.com/rubenbupe/recipe-video-parser
COPY . .
RUN make build

# Final image with shell and curl
FROM alpine:latest
WORKDIR /go
RUN apk add --no-cache curl
COPY --from=build /go/src/github.com/rubenbupe/recipe-video-parser/bin/api /go/bin/api
COPY --from=build /go/src/github.com/rubenbupe/recipe-video-parser/bin/cli /go/bin/cli

ENTRYPOINT ["/go/bin/api"]