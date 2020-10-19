This is a simple server that scrapes AWS RDS Instances

### Build
```
make
```

### Run
```
./aws_rds_exporter --aws.region=eu-west-1
```

## Exported Metrics

| Metric                              | Meaning                                                                                              | Labels                                        |
| ----------------------------------- | ---------------------------------------------------------------------------------------------------- | --------------------------------------------- |
| rds_storage   | Amount of storage in bytes for the RDS instance           | region, instance |

## Docker
You can deploy this exporter using the [alecrajeev/aws_rds_exporter](https://hub.docker.com/r/alecrajeev/aws_rds_exporter/) Docker Image.

Example
```
docker pull alecrajeev/aws_rds_exporter
docker run -p 9785:9785 alecrajeev/aws_rds_exporter
```

### Credentials
The `aws-rds_exporter` requires AWS credentials to access the AWS RDS API. For example you can pass them via env vars using `-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}` options.

