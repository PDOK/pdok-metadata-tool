FROM golang:1.24-alpine AS build-env

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags='-s -w' -o /pmt ./cmd/pmt

FROM scratch

ENV TZ=Europe/Amsterdam

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build-env /pmt /pmt

USER 65532:65532
WORKDIR /tmp

ENTRYPOINT ["/pmt"]
