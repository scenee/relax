#!/usr/bin/env bats

load test_helper

@test "parse Relfile" {
	cd ..

	source ../libexec/util-config

	eval $(_parse_yaml SampleApp/Relfile "config_")

	assert_success

	[ "$config_xcodeproj" = "SampleApp" ] \
	&& [ "$config_development_scheme" = "Sample App" ] \
	&& [ "$config_adhoc_export_options_thinning" = "iPhone7,1" ]

}
