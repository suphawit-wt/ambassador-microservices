# Ambassador Microservices

### Tools
- Framework: Fiber
- Database: MySQL, Redis
- ORM: GORM
- Message Queue: Kafka
- Payment: Stripe
- CI/CD: Docker, Kubernetes
- Authentication: JWT

## Development and Deployment

### Using Kafka Command in Docker Container

```bash
docker exec -it ambassador-kafka /bin/sh
```

```bash
cd /opt/bitnami/kafka/bin
```

List Topics

```bash
kafka-topics.sh --list --bootstrap-server 127.0.0.1:9092
```

Create Topic

```bash
kafka-topics.sh --create --bootstrap-server 127.0.0.1:9092 --topic mytopic --partitions 1 --replication-factor 1
```

Delete Topic

```bash
kafka-topics.sh --delete --topic mytopic --bootstrap-server 127.0.0.1:9092
```

Producer Console

```bash
kafka-console-producer.sh --bootstrap-server 127.0.0.1:9092 --topic mytopic
```

Consumer Console

```bash
kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic mytopic --from-beginning
```


### Create Kubernetes Secrets for Deployment on GKE

Create DB Secrets for Deployment

```bash
kubectl create secret generic db-secrets \
--from-literal=DB_HOST=<value> \
--from-literal=DB_USERNAME=<value> \
--from-literal=DB_PASSWORD=<value>
```

Create Redis Secrets for Deployment

```bash
kubectl create secret generic redis-secrets \
--from-literal=REDIS_HOST=<value>
```

Create Stripe Secrets for Deployment

```bash
kubectl create secret generic stripe-secrets \
--from-literal=STRIPE_KEY=<value>
```

Create SMTP Secrets for Deployment

```bash
kubectl create secret generic smtp-secrets \
--from-literal=SMTP_HOST=<value> \
--from-literal=SMTP_USERNAME=<value> \
--from-literal=SMTP_PASSWORD=<value>
```

Create Kafka Secrets for Deployment

```bash
kubectl create secret generic kafka-secrets \
--from-literal=KAFKA_SERVERS=<value> \
--from-literal=KAFKA_USERNAME=<value> \
--from-literal=KAFKA_PASSWORD=<value>
```

### Push Docker Image to Artifact Registry on GCP

Tag image to Artifact Registry on GCP

```bash
docker tag <image:tag> <gcp-region>-docker.pkg.dev/<project-id>/<repository>/<image:tag>
```

Push image to Artifact Registry on GCP

```bash
docker push <gcp-region>-docker.pkg.dev/<project-id>/<repository>/<image:tag>
```