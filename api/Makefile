build-MqttPublisherFunction:
	GOOS=linux CGO_ENABLE=0 go build -o lambda/mqttPublisher/bootstrap lambda/mqttPublisher/main.go
	cp lambda/mqttPublisher/bootstrap $(ARTIFACTS_DIR)/.


build-GetUserProfileImage:
	GOOS=linux CGO_ENABLE=0 go build -o lambda/getUserImageProfile/bootstrap lambda/getUserImageProfile/main.go
	cp lambda/getUserImageProfile/bootstrap $(ARTIFACTS_DIR)/.