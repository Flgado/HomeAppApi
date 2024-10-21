build-TriggerFunction:
	GOOS=linux CGO_ENABLE=0 go build -o lambda/addGroupScopeToIdToken/bootstrap lambda/addGroupScopeToIdToken/main.go
	cp lambda/addGroupScopeToIdToken/bootstrap $(ARTIFACTS_DIR)/.
