FROM golang:alpine

ADD . /go/src/github.com/alecrajeev/aws_rds_exporter
WORKDIR /go/src/github.com/alecrajeev/aws_rds_exporter

RUN go install

EXPOSE     9785
ENTRYPOINT [ "/go/bin/aws_rds_exporter" ]