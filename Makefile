.PHONY: test
test: bats
	@PATH="$$PWD/bats/bin:$$PWD/bin:$$PATH" test/run.sh

.PHONY: test-go
test-go:
	@PATH="$$PWD/bin:$$PATH" test/go.sh

version: 
	@PATH="$$PWD/bin:$$PATH" relax --version

bats:
	git clone --depth 1 https://github.com/sstephenson/bats.git

.PHONY: build
build:
	@export GOPATH=$$PWD/go; \
		>/dev/null pushd go; \
		go get -d ./src/...; \
		go install -a relparser; \
		go install -a lspp; \
		go install -a pputil; \
		popd > /dev/null;

.PHONY: man
man:
	@pushd share/man/man1 && ronn relax.1.ronn && popd
	@pushd share/man/man7 && ronn relax-*.7.ronn && popd
