# golang_messagio

## API
For more info refer to [API.md](API.md)
Check service status: `GET /healthcheck`
Send message: `POST /message`
Query message: `GET /query/message`
Query metrics: `GET /query/metrics`

## Deployment
```
git clone https://github.com/starnuik/golang_messagio
cd golang_messagio
DOCKER_BUILDKIT=1 docker-compose build
docker-compose up (-d)
```

## Components
* PostgreSQL - database
* Kafka - message broker
* pgmigrate - database migrations
* service-messagio - spec service
* service-fake-load - stress-test service, send 10e3-10e4 messages / 60 seconds