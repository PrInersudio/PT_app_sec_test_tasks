FROM golang:1.22

WORKDIR /usr/src/app

COPY task2/go.mod task2/go.sum ./
RUN go mod download && go mod verify

COPY task2/ config_for_docker.yaml ./
RUN go build -v -o /usr/local/bin/float_service.exe

EXPOSE 8081

CMD ["/usr/local/bin/float_service.exe", "config_for_docker.yaml"]