#!/usr/bin/bash

clear

PGOPATH="${GOPATH}"
export GOBIN=$GOPATH/bin

for dir in ./*; do
  [ -d "${dir}" ] || continue

  dirname="$(basename "${dir}")"
  [ "${dirname}" != "support" ] || continue

  absolute="$(realpath "${dir}")"
  GOPATH="${GOPATH}:${absolute}"

  pushd "${dirname}" > /dev/null || exit
  echo
  echo
  echo "Building ${dirname}..."
  go get
  go build "solution.go"

  time ./solution

  popd > /dev/null || exit
  GOPATH="${PGOPATH}"
done

#curl 'https://adventofcode.com/2020/day/8/input'
#-H 'User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0'
#-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8'
#-H 'Accept-Language: en-US,en;q=0.5'
#--compressed
#-H 'Referer: https://adventofcode.com/2020/day/8'
#-H 'Connection: keep-alive'
#-H 'Cookie: _ga=GA1.2.655323234.1607659979; _gid=GA1.2.1284925081.1607659979; session=53616c7465645f5f647f4538651604c3bb7e33b81900cbabe65dcea7dfa80ba782cc57cb189c7bb1b9f92cabc792fd53; _gat=1'
#-H 'Upgrade-Insecure-Requests: 1'
#-H 'Cache-Control: max-age=0'
#-H 'TE: Trailers'