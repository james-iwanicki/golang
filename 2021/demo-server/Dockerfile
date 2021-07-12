FROM golang:alpine

WORKDIR /build
COPY go.mod .
COPY . .
RUN go mod download
RUN go build demo-server.go

WORKDIR /dist
RUN cp /build/demo-server /dist/.

CMD [ "/build/demo-server" ]


