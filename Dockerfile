FROM golang:1.22-alpine as build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN  GOOS=linux GOARCH=amd64 go build -o bin/immudbapi cmd/*go

############################################### FINAL ##############################################
FROM alpine
ARG FILENAME="default.yaml"
ARG TARGETFILENAME="default.yaml"
ARG FILEDIR="./config"


WORKDIR /app
COPY --from=build /app/bin/immudbapi /app/bin/immudbapi
COPY $FILEDIR/${FILENAME} /app/config/${TARGETFILENAME}

#start command
ENTRYPOINT ["/app/bin/immudbapi"]




