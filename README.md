# Songs Search Exercise #

This repository implements an API to get song information from different services on Golang Echo.

## Install

Make sure you are using Docker

```
docker-compose up
```

## Routes

The service provides an unique endpoint with key authentication (default: token).

```
curl --location --request POST 'http://localhost:80/api/search' \
--header 'Authorization: token' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Ahora quién",
    "album": "Sigo Siendo Yo",
    "artist": "Marc Anthony"
}'
```

## Authentication

Use the Authorization header with a defined key.
You can change the key changing the AUTH_KEY env var with a `.env` file
