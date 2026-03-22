SUBPROJECTS := todo-cli guessr filesort snake memesweeper weathertape thumbforge dungeondice triviagoblin

.PHONY: deps-all build-all test-all cover-all clean-all list

deps-all:
	@set -e; \
	for p in $(SUBPROJECTS); do \
		echo "===> $$p: make deps"; \
		$(MAKE) -C $$p deps; \
	done

build-all:
	@set -e; \
	for p in $(SUBPROJECTS); do \
		echo "===> $$p: make build"; \
		$(MAKE) -C $$p build; \
	done

test-all:
	@set -e; \
	for p in $(SUBPROJECTS); do \
		echo "===> $$p: make test"; \
		$(MAKE) -C $$p test; \
	done

cover-all:
	@set -e; \
	for p in $(SUBPROJECTS); do \
		echo "===> $$p: make cover"; \
		$(MAKE) -C $$p cover; \
	done

clean-all:
	@set -e; \
	for p in $(SUBPROJECTS); do \
		echo "===> $$p: make clean"; \
		$(MAKE) -C $$p clean; \
	done

list:
	@echo $(SUBPROJECTS)
