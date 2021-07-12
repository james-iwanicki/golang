FROM golang:latest
WORKDIR /build

COPY go.mod .
COPY . .
RUN go mod download
RUN go build rest-api-server.go

WORKDIR /dist
RUN cp /build/rest-api-server /dist/.


CMD [ "/dist/rest-api-server" ]
