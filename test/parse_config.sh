#!/bin/bash

if test $# != 1; then
	echo "usage: $(basename $0) <yaml-file>"
	exit 1
fi

ME_DIR=$(dirname $0)

source ../libexec/util-config

# read yaml file

_parse_yaml $1 "config_" | awk '{
	if ( $0 ~ /.*\+=.*/ ) {
		match($0, /(.*)\+=(.*)/, a)
		print $0
		print "export "a[1]
	} else {
		print $0
	}
}' > temp-config

eval $(_parse_yaml $1 "config_")

destination=development

# access yaml content
eval echo '$'config_${destination}_configuration
eval echo '$'{config_${destination}_build_settings[@]}
eval echo '$'"{#config_${destination}_build_settings[@]}"

destination=adhoc
eval echo '$'config_${destination}_provisioning_profile

dests=( $(config_collect_releases $1 | tr '\r\n' ' ') )
echo ${dests[@]}
