FROM docker:dind
RUN apk add --no-cache go

ADD . /go/src/github.com/alecrajeev/aws_rds_exporter
WORKDIR /go/src/github.com/alecrajeev/aws_rds_exporter

RUN go install

USER nobody
EXPOSE     9785
ENTRYPOINT [ "/go/bin/aws_rds_exporter" ]