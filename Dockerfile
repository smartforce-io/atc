FROM golang:alpine as gobuilder
RUN apk add --update make git
ADD . /atc/
WORKDIR /atc/
RUN go mod download
RUN go build -o ../atc/bin/atcapp

FROM alpine

ARG GH_PEM_DATA

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
ENV ATC_APP_ID=79517
ENV ATC_PEM_DATA=${GH_PEM_DATA}

ENTRYPOINT ["/server/atcapp"]
