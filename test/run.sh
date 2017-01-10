#!/bin/bash

source PRIVATE # Defined 'CERTS_PASS'
export PATH="$PWD/bats/bin:$PWD/bin:$PATH"

relax keychain create relax.keychain -p relax
relax keychain add etc/RelaxCertificates.p12 -P "$CERTS_PASS"  -k relax.keychain -p relax
relax add etc/Relax_Development.mobileprovision
relax add etc/Relax_AdHoc.mobileprovision
export PROVISION_PROFILE_DEV="Relax Development"
export PROVISION_PROFILE_ADHOC="Relax AdHoc"

relax keychain use relax.keychain -p relax
bats test
relax keychain reset
