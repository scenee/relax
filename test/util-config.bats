#!/usr/bin/env bats

load test_helper

@test "config_load" {
	source $LIBEXEC_PATH/const
	source $LIBEXEC_PATH/util
	source $LIBEXEC_PATH/util-config

	export PATH=$LIBEXEC_PATH:$PATH
	make_temp
	run config_load Relfile adhoc

	[ "$REL_CONFIG_xcodeproj" = "SampleApp" ] \
	&& [ "$REL_CONFIG_adhoc_scheme" = "Sample App" ] \
	&& [ "$REL_CONFIG_adhoc_export_options_thinning" = "iPhone7,1" ]

	clean_temp
}

@test "config_get_distributions" {
	source $LIBEXEC_PATH/const
	source $LIBEXEC_PATH/util
	source $LIBEXEC_PATH/util-config

	export PATH=$LIBEXEC_PATH:$PATH
	run config_get_distributions Relfile

	[ "${#lines[@]}" = 7 ]
}


@test "relparser gen_source" {
	run $LIBEXEC_PATH/relparser -o $TMPDIR/source gen_source adhoc
	assert_success
	source $TMPDIR/source
	[ "$REL_CONFIG_xcodeproj" = "SampleApp" ] \
	&& [ "$REL_CONFIG_adhoc_scheme" = "Sample App" ] \
	&& [ "$REL_CONFIG_adhoc_export_options_thinning" = "iPhone7,1" ]
}
