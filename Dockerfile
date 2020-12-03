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

RUN mkdir /assets

USER atcuser

ENV ATC_PEM_PATH=/assets/github.pem
ENV ATC_CLIENT_ID=Iv1.afc8bdf21842ddc4
ENV ATC_APP_ID=79517

ENTRYPOINT ["/server/atcapp"]
