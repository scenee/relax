#!/usr/bin/env bats

load test_helper

@test "relparser source" {
	run $LIBEXEC_PATH/relparser -o $TMPDIR/source source adhoc
	assert_success
	source $TMPDIR/source
	[ "$REL_CONFIG_xcodeproj" = "SampleApp" ] \
	&& [ "$REL_CONFIG_adhoc_scheme" = "Sample App" ] \
	&& [ "$REL_CONFIG_adhoc_export_options_thinning" = "iPhone7,1" ]
}

@test "relparser export_options" {
	run $LIBEXEC_PATH/relparser -o $TMPDIR/options.plist -plist SampleApp/Info.plist export_options development
	assert_success
	cat $TMPDIR/options.plist
}
