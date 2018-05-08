/Schemes/ {
    while (1) {
        if ( (getline line) > 0 ) {
			if ( line ~ /^ *$/ ) {
				break
			} 
			gsub(/^  */,"", line)
			gsub(/ *$/,"", line)
            print line
        }
        else {
			break
        }
    }
}
