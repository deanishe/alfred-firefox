#!/usr/bin/env zsh

set -e

sizes=(19 38 48 96)
here="${${(%):-%x}:A:h}"
colour="c3c9c4"
font="elusive"
name="adjust"

for n in $sizes; do
	url="https://icons.deanishe.net/icon/${font}/${colour}/${name}/${n}.png"
	p="${here}/icons/icon-${n}.png"
	echo "fetching icon-${n}.png ..." >&2
	curl -L# "$url" > "$p"
done
