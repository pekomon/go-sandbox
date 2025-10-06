SUBPROJECTS := todo-cli

.PHONY: test-all
test-all:
	@set -e; \
	for p in $(SUBPROJECTS); do \
		echo "===> $$p: go test ./..."; \
		( cd $$p && go test ./... ); \
	done

.PHONY: list
list:
	@echo $(SUBPROJECTS)
