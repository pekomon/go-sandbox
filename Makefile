SUBPROJECTS := todo-cli guessr filesort snake brickbreaker memesweeper weathertape thumbforge dungeondice triviagoblin

.PHONY: deps-all build-all test-all cover-all clean-all clean-local-binaries list

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
	@$(MAKE) clean-local-binaries

clean-local-binaries:
	@set -e; \
	for p in $(SUBPROJECTS); do \
		echo "===> removing $$p/$${p##*/}"; \
		rm -f "$$p/$${p##*/}"; \
	done

list:
	@echo $(SUBPROJECTS)
