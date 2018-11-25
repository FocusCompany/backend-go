FROM postgres:10-alpine

RUN apk update && apk add git make musl-dev

## Note: the edge repositories often go down, if it's unreachable, just change dl-3 to use another number (e.g dl-4)
RUN apk add go --update-cache --repository http://dl-3.alpinelinux.org/alpine/edge/community --allow-untrusted

WORKDIR /go/src/github.com/FocusCompany/backend-go
ENV GOPATH=/go
RUN mkdir $GOPATH/bin
RUN wget -q -O - https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN go get -v -u -d github.com/golang-migrate/migrate/cli
RUN cd $GOPATH/src/github.com/golang-migrate/migrate/cli && \
	$GOPATH/bin/dep ensure && \
	go build -i -o $GOPATH/bin/migrate -tags "postgres"

# Custom initialization scripts, used to create DB and run migrations
COPY ./database/init-db.sh /docker-entrypoint-initdb.d/10-init-db.sh
COPY ./Makefile ./Makefile
COPY ./database/*.sql ./database/

RUN chmod 755 /docker-entrypoint-initdb.d/10-init-db.sh
