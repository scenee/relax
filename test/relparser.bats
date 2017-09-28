#!/usr/bin/env bats

load test_helper

@test "relparser source" {
	run $LIBEXEC_PATH/relparser -o $TMPDIR/source source adhoc
	assert_success
	source $TMPDIR/source
	[ "$REL_CONFIG_xcodeproj" = "SampleApp" ] \
	&& [ "$REL_CONFIG_adhoc_scheme" = "Sample App" ] \
	&& [ "$REL_CONFIG_adhoc_team_id" = "J3D7L9FHSS" ] \
	&& [ "$REL_CONFIG_adhoc_provisioning_profile" = "Relax AdHoc" ]
}

@test "relparser export_options" {
	run $LIBEXEC_PATH/relparser -o $TMPDIR/options.plist -plist SampleApp/Info.plist export_options development
	assert_success

	run plutil -lint $TMPDIR/options.plist
	assert_success

	cat $TMPDIR/options.plist > ../bats.log

	[ "$(/usr/libexec/PlistBuddy  -c "Print :teamID" $TMPDIR/options.plist)" == "J3D7L9FHSS" ] \
	&& [ "$(/usr/libexec/PlistBuddy  -c "Print :compileBitcode" $TMPDIR/options.plist)" == "false" ] \
	&& [ "$(/usr/libexec/PlistBuddy  -c "Print :uploadBitcode" $TMPDIR/options.plist)" == "true" ] 
}
