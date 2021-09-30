# Booker

A small CRUD application that manages books

## Running

1. Copy .env.example as .env, optionally change values to your taste
2. Run `docker-compose up -d`
3. Hit `localhost:8085/book` endpoint with POST, GET, PUT and DELETE methods

You can also run it locally with in-memory database. For this, set APP_PORT
and DB environment variables and run as follows:

```
go build cmd/booker/booker.go
DB=memory APP_PORT=8085 ./booker
```


Example of book payload for POST and PUT methods:

```
{
    "author": "Lewis Caroll",
    "title": "Alice in Wonderland",
    "publish_date": "2021-09-29T10:07:24.55176549+03:00",
    "publisher": "Publisher name",
    "rating": 3,
    "status": "checked_out"
}
```