#!/usr/bin/env bats

load test_helper

setup () {
	source libexec/const
	source libexec/util
	make_temp
}

teardown () {
	clean_temp
}


die_stdin () {
	echo "hello" | die
}

fin_stdin () {
	echo "hello" | fin
}

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

@test "die(stdin)" {
	run die_stdin
	assert_failure
	[[ "$output" =~ "hello" ]]
}

@test "fin(stdin)" {
	run fin_stdin 
	assert_success
	[[ "$output" =~ "hello" ]]
}
