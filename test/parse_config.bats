#!/usr/bin/env bats

load test_helper

@test "parse Relfile" {
	cd libexec

	source util-config

	cd ../test

	eval $(_parse_yaml SampleApp/Relfile "config_")

	assert_success

	[ "$config_xcodeproj" = "SampleApp" ]
	[ "$config_development_scheme" = "Sample App" ]
	[ "${#config_development_build_settings[@]}" = "2" ]
	[ "${config_development_build_settings[0]}" = "OTHER_SWIFT_FLAGS=-DMOCK" ]
	[ "$config_adhoc_export_options_thinning" = "iPhone7,1" ]
}
