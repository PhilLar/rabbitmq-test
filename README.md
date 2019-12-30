rabbitmq-test
=====

this is 2 servers: the producer subscribes to the matches channel via the websocket ("wss://ws-feed.pro.coinbase.com") and then publish message to the RabbitMQ "hello" queue. The consumer constantly read the queue and waits a message to log.


Installation and launch
------------
1) docker-compose up -d
2) go build -o bin/ ./cmd/...
3) ./bin/producer - run producer
   ./bin/consumer - run consumer