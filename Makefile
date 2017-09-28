.PHONY: build
build:
	@export GOPATH=$$PWD/go; \
		>/dev/null pushd go; \
		go get github.com/DHowett/go-plist; \
		go get gopkg.in/yaml.v2; \
		go install relparser; \
		go install lspp; \
		popd > /dev/null;

.PHONY: test
test: bats
	@PATH="$$PWD/bats/bin:$$PWD/bin:$$PATH" test/run.sh ${keychain}

version: 
	@PATH="$$PWD/bin:$$PATH" relax --version

bats:
	git clone --depth 1 https://github.com/sstephenson/bats.git

