FROM golang:1.10-alpine as build

RUN apk update && apk add --no-cache git protobuf-dev make zeromq zeromq-dev gcc  musl-dev
RUN wget -q -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get -u github.com/golang/protobuf/protoc-gen-go

WORKDIR /go/src/github.com/FocusCompany/backend-go

# Install dependencies
COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure --vendor-only

COPY . .

# Protobuf file generation
RUN make proto

RUN go install

#################
#    RUNTIME    #
#################
FROM alpine:latest as runtime

WORKDIR /run

COPY --from=build /go/bin/backend-go .
RUN apk update && apk add --no-cache zeromq

EXPOSE 8080
EXPOSE 5555

CMD [ "./backend-go"]
