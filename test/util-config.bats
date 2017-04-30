#!/usr/bin/env bats

load test_helper

@test "config_load" {
	source ../../libexec/const
	source ../../libexec/util
	source ../../libexec/util-config

	export PATH=../../libexec:$PATH
	make_temp
	run config_load Relfile adhoc

	[ "$REL_CONFIG_xcodeproj" = "SampleApp" ] \
	&& [ "$REL_CONFIG_adhoc_scheme" = "Sample App" ] \
	&& [ "$REL_CONFIG_adhoc_export_options_thinning" = "iPhone7,1" ]

	clean_temp
}

@test "config_get_distributions" {
	source ../../libexec/const
	source ../../libexec/util
	source ../../libexec/util-config

	export PATH=../../libexec:$PATH
	run config_get_distributions Relfile

	[ "${#lines[@]}" = 7 ]
}


@test "relparser gen_source" {
	run ../../libexec/relparser -o $TMPDIR/source gen_source adhoc
	assert_success
	source $TMPDIR/source
	[ "$REL_CONFIG_xcodeproj" = "SampleApp" ] \
	&& [ "$REL_CONFIG_adhoc_scheme" = "Sample App" ] \
	&& [ "$REL_CONFIG_adhoc_export_options_thinning" = "iPhone7,1" ]
}
