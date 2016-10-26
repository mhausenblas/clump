#!/bin/sh

# get the diagnostics bundle name:
bundle=$(dcos node diagnostics create all | tail -n 1 | rev | cut -f 1 -d ' ' | rev)
# wait a bit until the bundle is ready:
sleep 5
# download bundle to local machine:
dcos node diagnostics download $bundle
