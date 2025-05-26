FROM golang:alpine AS build

# Instala git y make
RUN apk add --update git make
WORKDIR /go/src/github.com/rubenbupe/recipe-video-parser
COPY . .
RUN make build

# Creates a minimal image with the application
FROM scratch
COPY --from=build /go/src/github.com/rubenbupe/recipe-video-parser/bin/api /go/bin/api
COPY --from=build /go/src/github.com/rubenbupe/recipe-video-parser/bin/cli /go/bin/cli
ENTRYPOINT ["/go/bin/api"]