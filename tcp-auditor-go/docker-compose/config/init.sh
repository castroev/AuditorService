#!/bin/sh

if [ ! -z "$MANIFESTCONTAINER" ]
then
    echo 'bootstrapping manifestservice'
    sh -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' http://${MANIFESTCONTAINER}/WebApplications)" != "200" ]]; do sleep 1; done'
    curl --request POST --data @/config/manifest.json http://${MANIFESTCONTAINER}/WebApplications/All -H "Content-Type:application/json"
fi

# If the container name has been set and publish app settings is true, then bootstrap that container
if [ ! -z "$CONSULCONTAINER" ] && [ $PUBLISHAPPSETTINGS = 1 ]; then
        for filename in /config/appsettings/*.json; do
                a="$(echo $filename | sed 's@.*/@@' | sed 's/[.].*//')"
                curl --request PUT --data @$filename http://${CONSULCONTAINER}/v1/kv/$a/appsettings.json -H "Content-Type:application/json"
        done
fi