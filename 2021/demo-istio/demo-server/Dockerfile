FROM golang:alpine

ENV version=development
WORKDIR /build
COPY go.mod .
COPY . .
RUN go mod download
RUN go build demo-server.go

WORKDIR /dist
RUN cp /build/demo-server /dist/.

CMD [ "/dist/demo-server" ]


