#!/usr/bin/env bats

load test_helper

setup () {
	source libexec/const
	source libexec/util
	source libexec/util-build
	make_temp
}

teardown () {
	clean_temp
}

@test "is_distribution_profile(Adhoc)" {
	run __is_distribution_profile sample/certificates/Relax_AdHoc.mobileprovision
	assert_success
	[ $output = "true" ]
}

@test "is_distribution_profile(Development)" {
	run __is_distribution_profile ./sample/certificates/Relax_Development.mobileprovision
	assert_success
	[ $output = "false" ]
}
	
@test "is_distribution_profile(Error)" {
	run __is_distribution_profile README.md
	assert_failure
}
