#!/bin/bash --login

# Xcode 7 (incl. 7.0.1) seems to have a dependency on the system ruby.
# xcodebuild is screwed up by using rvm to map to another non-system
# ruby. 
# This script's executed by /bin/bash --login

rvm use system
unset RUBYLIB
unset RUBYOPT
unset BUNDLE_BIN_PATH
unset _ORIGINAL_GEM_PATH
unset BUNDLE_GEMFILE

# Prevent a error 'shell_session_update: command not found'
shell_session_update() { :; }

set -ue

export PATH="$PWD/bats/bin:$PWD/bin:$PATH"
relax dec -k "$DECORD_KEY" etc/RelaxCertificates.p12.enc
relax dec -k "$DECORD_KEY" etc/Relax_Development.mobileprovision.enc
relax dec -k "$DECORD_KEY" etc/Relax_AdHoc.mobileprovision.enc
relax add -P "$CERTS_PASS" etc/RelaxCertificates.p12
relax add etc/Relax_Development.mobileprovision
relax add etc/Relax_AdHoc.mobileprovision

make test

relax rm "Relax Development"
relax rm "Relax AdHoc"
relax rm keychain
