# gokafka-exercise
Golang Kafka exercise simulates temperature monitoring. It is composed of a Kafka instance and two services, publisher and subscriber.  
The **publisher** service writes randomly generated temperature readings in Celsius every second to the `celcius-readings` Kafka topic.  
The **subscriber** service reads from the `celcius-readings` Kafka topic and output the data, a timestamp and the temperature, to stdout.  

## Prerequisites
- Install [Golang](https://golang.org/doc/install)
- Install [Docker](https://www.docker.com/get-started#h_installation) and [Docker Compose](https://docs.docker.com/compose/install/#install-compose)

## Running 
To run the services use `docker-compose up`  
The kafka service will start and output then the publisher and subscriber services will start. After a couple of seconds
the first messages received by the subscriber should output.  

Expected output:  
```
subscriber_1  | {"temperature": 85.46, "timestamp":1603762432}
subscriber_1  | {"temperature": 86.20, "timestamp":1603762433}
subscriber_1  | {"temperature": 33.54, "timestamp":1603762434}
subscriber_1  | {"temperature": 84.88, "timestamp":1603762435}
subscriber_1  | {"temperature": 50.77, "timestamp":1603762436}
subscriber_1  | {"temperature": 90.60, "timestamp":1603762437}
```

## Stopping
The service can be interrupted with `CTRL+C`  
If the Docker process does not clean up properly use `docker-compose stop`

## Running Externally
It might be useful to run the service allowing for external connections via the local machine.  
To do this use `docker-compose -f docker-compose-external.yml up`. Now, the service can be accessed from the local machine via `localhost:9093`  
Install [kafkacat](https://github.com/edenhill/kafkacat#install) and in a separate terminal window run `kafkacat -b localhost:9093 -C -t celcius-readings` to listen for messages.

## References
This service uses [bitnami-docker-kafka](https://github.com/bitnami/bitnami-docker-kafka)  
The kafka service in the docker-compose.yml is based off of [kafka-development-setup-example](https://github.com/bitnami/bitnami-docker-kafka#kafka-development-setup-example)  

