#!/usr/bin/ruby

require 'rubygems'
require 'yaml'
require 'json'

if ARGV.length != 3
	puts "usage: #{$0} <Relfile> <release> <Info.plist>"
	exit 1
end

relfile=ARGV[0]
release=ARGV[1]
plist_file=ARGV[2]


begin
	config = YAML.load_file(relfile)
rescue => e
	STDERR.puts e.message
	exit 1
end

plistbuddy="/usr/libexec/PlistBuddy"
INFO_PLIST_KEY="info_plist"

if config.has_key?(release) && config[release].has_key?(INFO_PLIST_KEY)
	config[release][INFO_PLIST_KEY].each do |k,v|
		case v

		when Array 
			puts "#{plistbuddy} -c \"Delete :#{k}\" #{plist_file}"
			%x[ #{plistbuddy} -c \"Delete :#{k}\" #{plist_file} ]
			puts "#{plistbuddy} -c \"Add :#{k} array\" #{plist_file}"
			%x[ #{plistbuddy} -c \"Add :#{k} array\" #{plist_file} ]

			v.each_index do |i|
				case v[i] 
				when String
					puts "#{plistbuddy} -c \"Add :#{k}:#{i} string #{v[i]}\" #{plist_file}"
					%x[ #{plistbuddy} -c \"Add :#{k}:#{i} string #{v[i]}\" #{plist_file} ]

				when Boolean
					puts "#{plistbuddy} -c \"Add :#{k}:#{i} bool #{v[i]}\" #{plist_file}"
					%x[ #{plistbuddy} -c \"Add :#{k}:#{i} bool #{v[i]}\" #{plist_file} ]

				else
					puts "Unsupported type #{v[i].class}: { #{k}: #{v[i]} }"
					exit 1
				end
			end

		when Hash
			puts "Unsupported type #{v.class}: { #{k}: #{v} }"

		else
			puts "#{plistbuddy} -c \"Set :#{k} #{v}\" #{plist_file}"
			%x[ #{plistbuddy} -c \"Set :#{k} #{v}\" #{plist_file} ]

		end
	end
end
