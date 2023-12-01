FROM golang:1.21.4-alpine3.18

WORKDIR /build
COPY go.mod .
COPY . .
RUN go mod download
RUN go build rest-api-client.go
WORKDIR /dist
RUN cp /build/rest-api-client /dist/.

CMD [ "/dist/rest-api-client" ]
