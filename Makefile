all: git2go
	go build

git2go:
	go get -d github.com/libgit2/git2go
	cd $$GOPATH/src/github.com/libgit2/git2go; \
		git submodule update --init; \
		make install

install: all
	find ~ -name prepare-commit-msg -exec cp prepare-commit-msg {} \;
