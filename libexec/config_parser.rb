#!/usr/bin/ruby
# TODO This script will replace _parse_yaml() in util-config

require 'rubygems'
require 'yaml'

if ARGV.length != 2
	puts "usage: #{$0} <Relfile> <release>"
	exit 1
end

relfile=ARGV[0]
release=ARGV[1]


begin
	config = YAML.load_file(relfile)
rescue => e
	STDERR.puts e.message
	exit 1
end

BUILD_SETTINGS_KEY="build_settings"

if config.has_key?(release) && config[release].has_key?(BUILD_SETTINGS_KEY)

	env_name="REL_CONFIG_#{release}_build_settings"

	puts "#{env_name}=()"

	config[release][BUILD_SETTINGS_KEY].each do |k,v|
		case v
		when String
			puts "#{env_name}+=(#{k}='#{v}')"
		when Array
			puts "#{env_name}+=(#{k}='#{v.join("{}")}')"
		else
			STDERR.puts "Found unsupported format in Relfile "
		end
	end

	puts "export #{env_name}"
end
