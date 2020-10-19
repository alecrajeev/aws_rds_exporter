ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="Alec Rajeev <alecinthecloud@gmail.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/aws_rds_exporter   /bin/aws_rds_exporter

USER nobody
EXPOSE     9785
ENTRYPOINT [ "/bin/aws_rds_exporter" ]