#!/bin/bash -eu

KEYCHAIN=relax.keychain

if [[ -f PRIVATE ]]; then
	source PRIVATE # Defined 'CERTS_PASS' values
fi

if [[ ${DECORD_KEY:-none} = "none" || ${CERTS_PASS:-none} = "none" ]]; then
	echo "Error: DECORD_KEY or CERTS_PASS are not defined."
	exit 1
fi

############################
# Set up a custom keychain #
############################

relax keychain create $KEYCHAIN -p relax

relax dec -p "$DECORD_KEY" sample/certificates/RelaxCertificates.p12.enc

relax keychain add sample/certificates/RelaxCertificates.p12 -P "$CERTS_PASS"  -k $KEYCHAIN -p relax

relax keychain use $KEYCHAIN -p relax

#################################
# Install Provisioning Profiles #
#################################

relax dec -p "$DECORD_KEY" sample/certificates/Relax_Development.mobileprovision.enc
relax dec -p "$DECORD_KEY" sample/certificates/Relax_AdHoc.mobileprovision.enc

relax profile add sample/certificates/Relax_Development.mobileprovision
relax profile add sample/certificates/Relax_AdHoc.mobileprovision

