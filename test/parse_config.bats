#!/usr/bin/env bats

load test_helper

@test "parse Relfile" {
	cd libexec
	source util-config
	cd ../test
	eval $(_parse_yaml Relfile-sample "config_")
	assert_success
	[ "$config_xcodeproj" = "SampleApp" ]
	[ "$config_SampleApp_sdk" = "iphoneos" ]
	[ "${#config_SampleApp_build_settings[@]}" = "2" ]
}
