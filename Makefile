.PHONY: build
build:
	@export GOPATH=$$PWD/go; \
		>/dev/null pushd go; \
		go get github.com/DHowett/go-plist; \
		go get gopkg.in/yaml.v2; \
		go install -a relparser; \
		go install -a lspp; \
		popd > /dev/null;

.PHONY: test
test: bats
	@PATH="/bin:$$PWD/bats/bin:$$PWD/bin:$$PATH" test/run.sh ${keychain}

.PHONY: version
version: 
	@PATH="$$PWD/bin:$$PATH" relax --version

.PHONY: bats
bats:
	git clone --depth 1 https://github.com/sstephenson/bats.git

.PHONY: man
man:
	@pushd share/man/man1 && ronn relax.1.ronn && popd
	@pushd share/man/man7 && ronn relax-*.7.ronn && popd
