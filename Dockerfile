FROM golang:1.17-alpine

WORKDIR /app

ADD . .

COPY go.mod .
COPY go.sum .
RUN go mod download

# COPY *.go ./
# # COPY /config/*.go ./config/.
# COPY .env ./

RUN go get owlint

RUN go build -o /owlint

EXPOSE 8080

CMD [ "/owlint" ]

# WORKDIR /app

# ADD . .

# RUN go get example

# RUN go install

# ENTRYPOINT [ "/api_go_test" ]