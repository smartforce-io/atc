FROM golang:alpine as gobuilder
RUN apk add --update make git
ADD . /atc/
WORKDIR /atc/
RUN go mod download
RUN go build -o ../atc/bin/atcapp

FROM alpine

RUN mkdir /atc
COPY --from=gobuilder /atc/bin/atcapp /atc/
WORKDIR /atc

RUN chmod +x atcapp
RUN adduser -D -g '' atcuser
RUN chown -R atcuser:atcuser /atc

USER atcuser

ENV GITHUB_TOKEN=$GITHUB_TOKEN
ENV FILE_TYPE=$FILE_TYPE
ENV COMMIT_SHA=$COMMIT_SHA
ENV GITHUB_REPOSITORY=$GITHUB_REPOSITORY
ENV BEHAVIOR = $BEHAVIOR
ENV TEMPLATE = $TEMPLATE
ENV REGEX = $REGEX
ENV CI_MODE = 'true'

ENTRYPOINT ["/atc/atcapp"]