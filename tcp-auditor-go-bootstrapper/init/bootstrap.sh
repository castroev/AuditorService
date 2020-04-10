#!/bin/sh

# If the container name has been set, then bootstrap that container
if [ -n "$CONSULCONTAINER" ] && [ "$PUBLISHENVSETTINGS" = 1 ];
then
    # To ensure we've given an istio sidecar time to bootstrap, we'll wait 2 seconds
    echo 'testing connection to consul'
    sh -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' http://${CONSULCONTAINER}/v1/kv/?keys)" != "200" ]]; do sleep 1; done'
    echo 'connection to consul successful'

    echo 'bootstrapping environment'
    curl -0 --request PUT --data @/config/appsettings.json http://"${CONSULCONTAINER}"/v1/kv/environment.json -H "Content-Type:application/json"
fi

if [ -n "$MANIFESTCONTAINER" ]
then
    echo 'bootstrapping manifestservice'
    sh -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' http://${MANIFESTCONTAINER}/WebApplications)" != "200" ]]; do sleep 1; done'
    curl -0 --request POST --data @/config/manifest.json http://"${MANIFESTCONTAINER}"/WebApplications/All -H "Content-Type:application/json"
fi
if [ -n "$CONSULCONTAINER" ] && [ "$PUBLISHAPPSETTINGS" = 1 ] && [ "$DEPLOYMENTSTAGE" = "CI" ]; then
    echo 'bootstrapping appsettings'
    for filename in /config/ci/appsettings/*.json; do
        echo "sending varfile $filename"
        a="$(echo "$filename" | sed 's@.*/@@' | sed 's/[.].*//')"
        curl -0 --request PUT --data @"$filename" http://"${CONSULCONTAINER}"/v1/kv/"$a"/appsettings.json -H "Content-Type:application/json"
    done
fi

if [ -n "$CONSULCONTAINER" ] && [ "$PUBLISHAPPSETTINGS" = 1 ] && [ "$DEPLOYMENTSTAGE" = "QA" ]; then
    echo 'bootstrapping appsettings'
    for filename in /config/qa/appsettings/*.json; do
        a="$(echo "$filename" | sed 's@.*/@@' | sed 's/[.].*//')"
        curl -0 --request PUT --data @"$filename" http://"${CONSULCONTAINER}"/v1/kv/"$a"/appsettings.json -H "Content-Type:application/json"
    done
fi

if [ -n "$CONSULCONTAINER" ] && [ "$PUBLISHAPPSETTINGS" = 1 ] && [ "$DEPLOYMENTSTAGE" = "PROD" ]; then
    echo 'bootstrapping appsettings'
    for filename in /config/prod/appsettings/*.json; do
        a="$(echo "$filename" | sed 's@.*/@@' | sed 's/[.].*//')"
        curl -0 --request PUT --data @"$filename" http://"${CONSULCONTAINER}"/v1/kv/"$a"/appsettings.json -H "Content-Type:application/json"
    done
fi
