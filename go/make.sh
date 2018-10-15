#!/usr/bin/env bash

if [ $1 ]
then
    if [ $1 == "init" ]
    then
        go build src/init/init.go src/init/database.go src/init/logger.go
    elif [ $1 == "elk" ]
    then
        go build src/elk/elk.go src/elk/database.go src/elk/elkClient.go src/elk/indexer.go src/elk/logger.go
    elif [ $1 == "api" ]
    then
        go build src/api/api.go src/api/database.go src/api/model.go src/api/logger.go
    else
        echo "No valid arguments found"
    fi
fi