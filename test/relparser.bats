#!/usr/bin/env bats

load test_helper
source libexec/util-build

@test "relparser source" {
	run $LIBEXEC_PATH/relparser source adhoc
	assert_success
	source "${output}"
	[ "$REL_CONFIG_xcodeproj" = "SampleApp" ] \
		&& [ "$_SCHEME" = "Sample App" ] \
		&& [ "$_TEAM_ID" = "J3D7L9FHSS" ] \
		&& [ "$_PROVISIONING_PROFILE" = "Relax AdHoc" ]
}

@test "relparser export_options" {
	run $LIBEXEC_PATH/relparser -o $TMPDIR/options.plist -plist SampleApp/Info.plist export_options development
	assert_success

	run plutil -lint $TMPDIR/options.plist
	assert_success

	# cat $TMPDIR/options.plist > ../bats.log

	[ "$(/usr/libexec/PlistBuddy  -c "Print :teamID" $TMPDIR/options.plist)" == "J3D7L9FHSS" ] \
	&& [ "$(/usr/libexec/PlistBuddy  -c "Print :compileBitcode" $TMPDIR/options.plist)" == "false" ] \
	&& [ "$(/usr/libexec/PlistBuddy  -c "Print :uploadBitcode" $TMPDIR/options.plist)" == "false" ] 
}

@test "relparser plist" {
	run $LIBEXEC_PATH/relparser -o $TMPDIR/new.plist -plist SampleApp/Info.plist plist development
	cat $TMPDIR/new.plist
	run plutil -lint $TMPDIR/new.plist
	assert_success
}

