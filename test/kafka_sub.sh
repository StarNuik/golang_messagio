sudo docker exec -it golang_messagio_kafka_1 \
    /opt/bitnami/kafka/bin/kafka-console-consumer.sh \
    --bootstrap-server localhost:9092 \
    --from-beginning \
    --topic $1