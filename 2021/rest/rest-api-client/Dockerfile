FROM golang:1
WORKDIR /build

COPY go.mod .
COPY . .
RUN go mod download
RUN go build rest-api-client.go

WORKDIR /dist
RUN cp /build/rest-api-client /dist/.

#ENV GODEBUG="x509ignoreCN=0"
CMD [ "/dist/rest-api-client" ]
#CMD [ "/bin/bash", "-c", "GODEBUG=x509ignoreCN=0 && /dist/rest-api-client" ]
