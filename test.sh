#!/bin/bash

cwd="$(pwd)"

testFile="$cwd/testFile$(date +'%s')"
testEnc="$cwd/testEncFile$(date +'%s')"
testDec="$cwd/testDecFile$(date +'%s')"
HBIN="$cwd/huffandpuff"

cat << EOF > "$testFile"
This is a test file blah blah blah
blergh blergh bleueueueuee
EOF

"$HBIN" -c -in "$testFile" -out "$testEnc"
"$HBIN" -d -in "$testEnc" -out "$testDec"

if ! cmp "$testFile" "$testDec"; then
  echo "mismatch"
fi

rm "$testFile" "$testEnc" "$testDec"