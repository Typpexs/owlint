FROM golang:1.17-alpine

WORKDIR /app

ADD . .

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN go get owlint

RUN go build -o /owlint

EXPOSE 8080

CMD [ "/owlint" ]