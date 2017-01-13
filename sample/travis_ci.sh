#!/bin/bash --login

# Xcode 8 or later seems to have a dependency on the system ruby.
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

export PATH="$PWD/bin:$PATH"

#########################
# Set up relax keychain #
#########################

relax keychain create relax.keychain -p relax
relax dec -p "$DECORD_KEY" sample/certificates/RelaxCertificates.p12.enc
relax keychain add sample/certificates/RelaxCertificates.p12 -P "$CERTS_PASS"  -k relax.keychain -p relax

############################
# Install mobileprovisions #
############################

relax dec -p "$DECORD_KEY" sample/certificates/Relax_Development.mobileprovision.enc
relax dec -p "$DECORD_KEY" sample/certificates/Relax_AdHoc.mobileprovision.enc
relax add sample/certificates/Relax_Development.mobileprovision
relax add sample/certificates/Relax_AdHoc.mobileprovision
export PROVISION_PROFILE_DEV="Relax Development"
export PROVISION_PROFILE_ADHOC="Relax AdHoc"

###########
# Do Test #
###########

relax keychain use relax.keychain -p relax
make test
relax keychain reset

############
# Teardown #
############

relax rm "Relax Development"
relax rm "Relax AdHoc"

relax keychain delete relax.keychain
