This is a simple server that scrapes AWS RDS Instances

### Build
```
make
```

### Run
```
./aws_rds_exporter --aws_rds.region=us-east-1
```

## Exported Metrics

| Metric                              | Meaning                                                                                              | Labels                                        |
| ----------------------------------- | ---------------------------------------------------------------------------------------------------- | --------------------------------------------- |
| aws_rds_storage   | Amount of storage in bytes for the RDS instance           | region, instance |
| aws_rds_iops   | Amount of iops for the RDS instance           | region, instance |

### Flags

```bash
./aws_rds_exporter --help
```

* __`aws_rds.region`:__ AWS Region to run API calls against.

## Unit Tests
Use the below to run unit tests locally.
```
go list ./... | grep -v /vendor/ | go test -v
```

## Docker
You can deploy this exporter using the [alecrajeev/aws_rds_exporter](https://hub.docker.com/r/alecrajeev/aws_rds_exporter/) Docker Image.

Example
```
docker pull alecrajeev/aws_rds_exporter
docker run -p 9785:9785 alecrajeev/aws_rds_exporter
```

### Credentials
The `aws_rds_exporter` requires AWS credentials to access the AWS RDS API. For example you can pass them via env vars using `-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}` options.

