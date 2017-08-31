.PHONY: build
build:
	@cd go; \
		export GOPATH=$$PWD; \
		go get github.com/DHowett/go-plist; \
		go get gopkg.in/yaml.v2; \
		pushd src/relparser >/dev/null; \
		go install; \
		popd > /dev/null; \
		pushd src/lspp >/dev/null; \
		go install; \
		popd > /dev/null;

.PHONY: test
test: bats
	@PATH="$$PWD/bats/bin:$$PWD/bin:$$PATH" test/run.sh ${keychain}

version: 
	@PATH="$$PWD/bin:$$PATH" relax --version

bats:
	git clone --depth 1 https://github.com/sstephenson/bats.git

