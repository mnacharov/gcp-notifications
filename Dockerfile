FROM golang:1.19.5-bullseye

WORKDIR /usr/src/gcp-notifications

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/gcp-notifications ./...

CMD ["/usr/local/bin/gcp-notifications"]
