#!/bin/sh

set -e

TAG=$(git tag | tail -1)
GOVER=$(go version | grep -oE 'go[0-9]\.[0-9]+')

PREV="testdata/bench_${TAG}_${GOVER}"
LATEST="testdata/bench_@latest_$GOVER"

if [ ! -f "$PREV" ]; then
    git mv "$LATEST" "$PREV"
    go test -run '^$' -bench . -count 5 > "$LATEST"
fi

benchstat "$PREV" "$LATEST"
