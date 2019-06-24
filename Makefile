all:
	go build

install: all
	find ~ -name prepare-commit-msg -exec cp prepare-commit-msg {} \;

.PHONY: all install
