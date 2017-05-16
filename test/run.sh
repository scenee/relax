#!/bin/bash -eu

if [[ $# = 0 ]]; then
	keychain=relax.keychain
else
	keychain=$keychain
	if ! grep -q "\.keychain" <<< "$keychain"; then
		echo "Error: invalid keychain name($1)"
		echo "    Please add '.keychain' extension."
		exit 1
	fi
fi

if [[ -f PRIVATE ]]; then
	source PRIVATE # Defined 'CERTS_PASS' values
fi

if [[ ${DECORD_KEY:-none} = "none" || ${CERTS_PASS:-none} = "none" ]]; then
	echo "Error: DECORD_KEY or CERTS_PASS are not defined."
	exit 1
fi

###################
# Set up keychain #
###################

relax keychain create $keychain -p relax

relax dec -p "$DECORD_KEY" sample/certificates/RelaxCertificates.p12.enc

relax keychain add sample/certificates/RelaxCertificates.p12 -P "$CERTS_PASS"  -k $keychain -p relax

############################
# Install mobileprovisions #
############################

relax dec -p "$DECORD_KEY" sample/certificates/Relax_Development.mobileprovision.enc
relax dec -p "$DECORD_KEY" sample/certificates/Relax_AdHoc.mobileprovision.enc

relax profile add sample/certificates/Relax_Development.mobileprovision
relax profile add sample/certificates/Relax_AdHoc.mobileprovision

export PROVISION_PROFILE_DEV="Relax Development"
export PROVISION_PROFILE_ADHOC="Relax AdHoc"

############
# Run Test #
############


trap "relax keychain reset" EXIT
relax keychain use $keychain -p relax

export NOCOLOR=true
bats test

######################
# Tear down keychain #
######################

relax profile rm "Relax Development"
relax profile rm "Relax AdHoc"

relax keychain delete $keychain
