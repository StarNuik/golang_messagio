sudo docker exec -it golang_messagio_kafka_1 \
    /opt/bitnami/kafka/bin/kafka-console-producer.sh \
    --bootstrap-server localhost:9092 \
    --topic $1