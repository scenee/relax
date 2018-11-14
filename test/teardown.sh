#!/bin/bash -eu

#############
# Tear down #
#############

KEYCHAIN=relax.keychain

if [[ -f PRIVATE ]]; then
	source PRIVATE # Defined 'CERTS_PASS' values
fi

relax profile rm "Relax Development"
relax profile rm "Relax AdHoc"

# Test `relax keychain rm`
while IFS= read -r identity && [[ -n "$identity" ]]
do
	relax keychain rm "$(echo "$identity" | awk '{ print $1 }')" -k $KEYCHAIN -p relax
done < <(relax keychain info sample/certificates/RelaxCertificates.p12 -P "$CERTS_PASS")

relax keychain delete $KEYCHAIN
