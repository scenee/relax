#!/bin/bash -eu

trap "relax keychain reset" EXIT

test/setup.sh

############
# Run Test #
############

export NOCOLOR=true
bats test

test/teardown.sh
