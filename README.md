# golang_messagio

## API
For more info refer to [API.md](API.md)<br>
Check service status: `GET /healthcheck`<br>
Send message: `POST /message`<br>
Query message: `GET /query/message`<br>
Query metrics: `GET /query/metrics`<br>

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