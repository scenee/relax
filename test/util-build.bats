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
	run __is_distribution_profile sample/certificates/Relax_Development.mobileprovision
	assert_success
	[ $output = "false" ]
}
	
@test "is_distribution_profile(Error)" {
	run __is_distribution_profile README.md
	assert_failure
}

@test "dec_provisioning_profile" {
	run dec_provisioning_profile sample/certificates/Relax_Development.mobileprovision -o $REL_TEMP_DIR/Relax_Development.mobileprovision.dec
	assert_success
	run grep -q "Relax Development" $REL_TEMP_DIR/Relax_Development.mobileprovision.dec
	assert_success
}

@test "dec_provisioning_profile(Error)" {
	run dec_provisioning_profile sample/certificates/NotFound.mobileprovision -o $REL_TEMP_DIR/Relax_Development.mobileprovision.dec
	assert_failure
}

@test "find_mobileprovision(Not found)" {
	run find_mobileprovision "Not Found"
	assert_success
	[ -z "$output" ]
}

@test "find_mobileprovision(Relax AdHoc)" {
	run find_mobileprovision "Relax AdHoc"
	assert_success
	[ -n "$output" ]
}
#relax profile add sample/certificates/Relax_AdHoc.mobileprovision
