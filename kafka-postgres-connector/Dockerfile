FROM gcr.io/google_samples/k8skafka:v1

RUN wget -qO - https://packages.confluent.io/deb/4.0/archive.key | apt-key add - && \
    apt update && apt install -y software-properties-common apt-transport-https && \
    add-apt-repository "deb [arch=amd64] https://packages.confluent.io/deb/4.0 stable main" && \
    apt update && apt install -y confluent-kafka-connect-jdbc && \
    apt remove -y software-properties-common apt-transport-https

WORKDIR /opt/kafka/bin
COPY jdbc-sink.properties .
COPY worker.properties .
ENV CLASSPATH /usr/share/java/kafka-connect-jdbc/*
CMD sh connect-standalone.sh worker.properties jdbc-sink.properties
