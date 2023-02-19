This project consists of two services.

The **READER** service reads a json file with data (or generates a file with fake data) and serves two handlers.
- /generate - Generates a file ex:`http://localhost:8002/generate?length=500`
- /stats - Gets parameters and returns a representation of the aggregated data ex:`http://localhost:8002/stats?interval=year&start=1595575638&end=1637685638`

Reads a file at startup using flags. ex:`$ go run . -r=data.json -g -l=10000`
- -r=filename - read file
- -g  (bool)  - generate fake data file
- -l=100 - length of fake data

The stream data is sent via RPC to the second(storage) service.

The **STORAGE** service is a RPC server, it receives data and requests from the first service, processes them and interacts with the database - elastic.

Run with docker-compose 
-

`docker-compose up --build`

`docker-compose up --build -d`

`docker-compose up -d`

`docker-compose down`

Or separately each service and elastic by editing docker-compose.local.yml

`docker-compose -f docker-compose.local.yml up -d`
