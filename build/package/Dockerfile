FROM golang:1.20 AS builder
LABEL stage=intermediate
ARG VERSION
COPY . /backend
WORKDIR /backend
ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build -ldflags "-X main.commitHash=$(git rev-parse HEAD)" -o backend cmd/backend/*

FROM alpine:latest
LABEL maintainer="haveachin <haveachin@haveachin.de>"
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /backend/backend ./
RUN chmod +x ./backend
ENTRYPOINT [ "./backend" ]