_relax() {
	local cur prev words cword

	COMPREPLY=()
	cur=${COMP_WORDS[COMP_CWORD]}
	prev=${COMP_WORDS[COMP_CWORD-1]}
	words=("${COMP_WORDS[@]}")

	local module special i
	for (( i=0; i < ${#words[@]}-1; i++ )); do
		commands="$(relax commands | grep -v -- --version | tr " " "\|")"
		if echo ${words[i]} | grep -q "$commands" ; then
		    special=${words[i]}
		fi
		modules="$(relax commands --modules | tr " " "\|" )"
		if echo ${words[i]} | grep -q "$modules" ; then
		    module=${words[i]}
		fi
	done

	if [[ -n $module ]]; then
		case $module in
		keychain)
			if [[ $prev = $module ]]; then
				COMPREPLY=( $( compgen -W "$(relax $module completions)" -- $cur ) )
			else
				command=${words[2]}
				case $command in
				reset)
					;;
				*)
					COMPREPLY=( $( compgen -W "$(relax $module completions ${command} $prev $cur)" -- $cur ) )
					;;
				esac
			fi
			return
			;;
		esac
	fi

	if [[ -n $special ]]; then
		case $special in
			archive|build|export|show|package)
				if [[ $prev = $special ]]; then
					COMPREPLY=( $( compgen -W "$(relax completions releases)" -- $cur ) )
				else
					COMPREPLY=( $( compgen -W "$(relax $special --complete $prev)" -- $cur ) )
				fi
				return
				;;
			upload)
				if [[ $prev = $special ]]; then
					COMPREPLY=( $( compgen -W "$(relax $special --complete)" -- $cur ) )
				fi
				return
				;;
			*)
				_filedir
				return
				;;
		esac
	fi

	case "$prev" in
		--config)
			_filedir
			;;
	esac

	case "$cur" in
		*)
		COMPREPLY=( $(compgen -W "$(relax commands) $(relax commands --modules)" -- "$cur") )
		;;
	esac
}

complete -F _relax  relax
