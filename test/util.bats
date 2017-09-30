#!/usr/bin/env bats

load test_helper

setup () {
	source libexec/const
	source libexec/util
	make_temp
	stdout_path="${BATS_TMPDIR}/bats.stdout"
	stderr_path="${BATS_TMPDIR}/bats.stderr"
}

teardown () {
	clean_temp
}


#die_stdin () {
#	echo "hello" | die
#}
#
#fin_stdin () {
#	echo "hello" | fin
#}

@test "die" {
	run die "hello"
	assert_failure
	[[ "$output" =~ "hello" ]]
}

@test "fin" {
	run fin "hello"
	assert_success
	[[ "$output" =~ "hello" ]]
}

@test "loge" {

	loge "hello" 2>${stderr_path}

	grep ".*Error:.*hello"  ${stderr_path}

	loge "hello" 1>${stdout_path}

	! grep ".*Error:.*hello"  ${stdout_path}
}

@test "logw" {

	logw "hello" 2>${stderr_path}

	grep ".*WARNING:.*hello"  ${stderr_path}

	logw "hello" 1>${stdout_path}

	! grep ".*WARNING:.*hello"  ${stdout_path}
}

#@test "die(stdin)" {
#	run die_stdin
#	assert_failure
#	[[ "$output" =~ "hello" ]]
#}
#
#@test "fin(stdin)" {
#	run fin_stdin 
#	assert_success
#	[[ "$output" =~ "hello" ]]
#}
