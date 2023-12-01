FROM golang:1.21.4-alpine3.18

ENV version=development
WORKDIR /build
COPY go.mod .
COPY . .
RUN go mod download
RUN go build demo-server.go

WORKDIR /dist
RUN cp /build/demo-server /dist/.

CMD [ "/dist/demo-server" ]


