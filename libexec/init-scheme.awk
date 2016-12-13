/Schemes/ {
	print "xcode_schemes=()"
    while (1) {
        if ( (getline foo) > 0 ) {
			if ( foo ~ /^ *$/ ) {
				break
			} 
			gsub(/^  */,"", foo)
			gsub(/ *$/,"", foo)
            print "xcode_schemes+=(\""foo"\")"
        }
        else {
			break
        }
    }
	print "export xcode_schemes"
}
