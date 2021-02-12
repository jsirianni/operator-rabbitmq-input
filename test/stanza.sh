#!/bin/bash

set -e

cd "$(dirname "$0")"

mkdir -p stanza &&

wget https://github.com/observIQ/stanza/archive/v0.13.12.tar.gz
tar -xvf v0.13.12.tar.gz --strip-components 1 -C stanza/
rm -f v0.13.12.tar.gz

cp go.mod stanza/cmd/stanza/go.mod
cp init_common.go stanza/cmd/stanza/init_common.go

cd stanza/cmd/stanza
go build -o ../../../stanza-test
cd ../../../

./stanza-test
