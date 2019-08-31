FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD /src /src
RUN cd /src && go test && go build -o goapp

FROM alpine
WORKDIR /app
COPY ./ipsum.txt /tmp
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp