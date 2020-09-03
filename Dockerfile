FROM golang:alpine as gobuilder
RUN apk add --update make
ADD . /atc/
WORKDIR /atc/
ENV GOPATH=$GOPATH:/atc
RUN make all

FROM alpine
EXPOSE 8080

RUN mkdir /server
COPY --from=gobuilder /atc/bin/atcapp /server/
WORKDIR /server

RUN chmod +x atcapp
RUN adduser -D -g '' atcuser
RUN chown -R atcuser:atcuser /server

USER atcuser

ENTRYPOINT ["/server/atcapp"]
