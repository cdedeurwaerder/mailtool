FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o mailtool ./cmd/main.go

#===========================

FROM alpine:latest

RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite

COPY --from=builder /src/mailtool /bin/mailtool

RUN mkdir /data
ENV MAILTOOL_DATADIR=/data

CMD ["/bin/mailtool"]