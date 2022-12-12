# modules caching
FROM golang:1.19-bullseye as modules

COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# builder
FROM golang:1.19-bullseye as build

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -ldflags="-s -w" -o /go/bin/linebot ./cmd/app

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/linebot /
# ENTRYPOINT [ "/linebot" ]
CMD ["/linebot"]

