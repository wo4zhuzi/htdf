# written by junying, 2020-01-15

genesis2accs:
	@printkey accounts genesis_old.json address amount denom|repl "[u'" ""|repl "']" ""|deline stake|rowcnt 2 > accounts.list

checksum:
	@printkey accounts genesis_old.json address amount denom|repl "[u'" ""|repl "']" ""|deline stake|rowcnt 2|column 2|sumup