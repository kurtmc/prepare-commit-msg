all: git2go
	go build

git2go:
	go get -d gopkg.in/libgit2/git2go.v26
	cd ../../../gopkg.in/libgit2/git2go.v26/; \
		git submodule update --init; \
		make install

install: all
	find ~ -name prepare-commit-msg -exec cp prepare-commit-msg {} \;
