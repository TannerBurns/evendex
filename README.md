# EvenDex

# How to use

    # Go

        inside 'src' folder

        go build views.go model.go indexer.go

        ./views /path/to/config/config.conf

        go api hosted: http://localhost:8000/api/v1/events
        
        #Example Config.conf:
        
            [postgresql]
            host=localhost
            database=evendex
            user=postgres
            password=password
    
    # Flask

        inside 'flask' folder

        python views.py

        flask frontend hosted: http://localhost:5000/events

