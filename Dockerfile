FROM node:8 as js

ADD . /app
WORKDIR /app/assets

RUN npm i
RUN npm run build

FROM golang:1.9 as compiler

COPY ./ /go/src/github.com/mikeifomin/graphql-subscriptions-go-example
WORKDIR  /go/src/github.com/mikeifomin/graphql-subscriptions-go-example
COPY --from=js /app/assets/dist assets/dist

RUN go get -u github.com/jteeuwen/go-bindata/...

RUN cd assets && go-bindata -prefix dist -o bindata.go -pkg assets dist/... 
RUN go build -o /bin/gqlexample main.go 

FROM debian:jessie
COPY --from=compiler /bin/gqlexample /bin/gqlexample
EXPOSE 80 443
ENTRYPOINT ["/bin/gqlexample"]
