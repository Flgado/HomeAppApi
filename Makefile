build-MqttPublisherFunction:
	GOOS=linux CGO_ENABLE=0 go build -o lambda/mqttPublisher/bootstrap lambda/mqttPublisher/main.go
	cp lambda/mqttPublisher/bootstrap $(ARTIFACTS_DIR)/.


build-Test:
	GOOS=linux CGO_ENABLE=0 go build -o lambda/test/bootstrap lambda/test/main.go
	cp lambda/test/bootstrap $(ARTIFACTS_DIR)/.