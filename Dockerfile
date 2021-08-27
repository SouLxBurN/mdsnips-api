FROM golang:1.16.5-alpine

WORKDIR /src
COPY ./ ./

RUN go install

CMD "soulxsnips"

EXPOSE 3000
