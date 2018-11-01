#!/usr/bin/env bash

if [ $1 ]
then
    if [ $1 == "init" ]
    then
        echo "Building init.."
        go build src/api/init.go src/api/database.go src/api/logger.go src/api/elkClient.go
    elif [ $1 == "api" ]
    then
        echo "Building api.."
        go build src/api/api.go src/api/database.go src/api/model.go src/api/logger.go src/api/elkClient.go src/api/esmodel.go
    elif [ $1 == "all" ]
    then
        echo "Building all"
        echo "########################"
        echo "Building init.."
        go build src/api/init.go src/api/database.go src/api/logger.go src/api/elkClient.go
        echo "Building api.."
        go build src/api/api.go src/api/database.go src/api/model.go src/api/logger.go src/api/elkClient.go src/api/esmodel.go
    elif [ $1 == "build" ]
    then
        echo "Building init.."
        go build src/api/init.go src/api/database.go src/api/logger.go src/api/elkClient.go
        ./init conf/config.dev.conf
        rm init
        echo "Building api.."
        go build src/api/api.go src/api/database.go src/api/model.go src/api/logger.go src/api/elkClient.go src/api/esmodel.go
    fi
    echo "Done!"
fi
