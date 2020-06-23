# message-api

Technical test

CURL examples, for testing core endpoints:
-  curl -v -H "Accept: application/json" -H "Content-t
  ype: application/json" -X POST -d ' {"name":"br", "payload":"payload test", "metadata":{"timestamp":"2020-12-12 19:02:33.123", "uui
  d":"10eec78c-b4a3-11ea-b3de-0242ac130004"}}'  http://localhost:9090/messages
- curl -X  GET  http://localhost:9090/messages
- curl -X  GET  http://localhost:9090/messages/bbr


To run in Docker execute command:
- docker-compose up --build

To run tests in Docker execute:
- docker-compose  -f docker-compose.test.yml up  --build --abort-on-container-exit
