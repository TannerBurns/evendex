#!/usr/bin/env bash

if [ $1 ]
then
    if [ $1 == "init" ]
    then
        echo "Building init.."
        go build src/init/init.go src/init/database.go src/init/logger.go
    elif [ $1 == "elk" ]
    then
        echo "Building elk.."
        go build src/elk/elk.go src/elk/database.go src/elk/elkClient.go src/elk/indexer.go src/elk/logger.go
    elif [ $1 == "api" ]
    then
        echo "Building api.."
        go build src/api/api.go src/api/database.go src/api/model.go src/api/logger.go
    elif [ $1 == "all" ]
        echo "Building all"
        echo "########################"
        echo "Building init.."
        go build src/init/init.go src/init/database.go src/init/logger.go
        echo "Building elk.."
        go build src/elk/elk.go src/elk/database.go src/elk/elkClient.go src/elk/indexer.go src/elk/logger.go
        echo "Building api.."
        go build src/api/api.go src/api/database.go src/api/model.go src/api/logger.go
    fi
    echo "Done!"
else
    echo "Building all"
    echo "########################"
    echo "Building init.."
    go build src/init/init.go src/init/database.go src/init/logger.go
    echo "Building elk.."
    go build src/elk/elk.go src/elk/database.go src/elk/elkClient.go src/elk/indexer.go src/elk/logger.go
    echo "Building api.."
    go build src/api/api.go src/api/database.go src/api/model.go src/api/logger.go
    echo "Done!"
fi