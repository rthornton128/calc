#!/bin/bash

# Change if you want it to go somewhere else
prefix="$GOPATH/bin"

if [ "$1" = "clean" ]; then
	cd calcc;
	go clean -i;
	cd "..";
	rm -f "$prefix/calc.bash";
	exit 0;
fi

echo "Running test suite..."

# BUG If the final test passes, then go test has a non-zero exit status
go test ./...

if [ "$?" -eq 0 ]; then
	echo "Installing calcc to $prefix..."
	cd calcc;
	go install;
	cd "..";
	echo "Installing calc script to $prefix..."
	install "calc.bash" "$prefix";
fi
echo "Finished"
