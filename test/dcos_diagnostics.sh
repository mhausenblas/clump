#!/bin/sh

bundle=$(dcos node diagnostics create all | tail -n 1 | rev | cut -f 1 -d ' ' | rev)
sleep 10
dcos node diagnostics download $bundle
