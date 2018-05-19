#!/bin/bash --login

# Xcode 8 or later seems to have a dependency on the system ruby.
# xcodebuild is screwed up by using rvm to map to another non-system
# ruby. 
# This script's executed by /bin/bash --login

rvm use system
unset RUBYLIB
unset RUBYOPT
unset BUNDLE_BIN_PATH
unset _ORIGINAL_GEM_PATH
unset BUNDLE_GEMFILE

# Prevent a error 'shell_session_update: command not found'
shell_session_update() { :; }

set -ue

make test
