FROM golang:1.19.5-bullseye as build

WORKDIR /go/src/gcp-notifications
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -v -o /go/bin/gcp-notifications ./...


FROM gcr.io/distroless/static-debian11:latest

COPY --from=build /go/bin/gcp-notifications /go/src/gcp-notifications/slack.json /

CMD ["/gcp-notifications"]
