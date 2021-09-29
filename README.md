# Booker

A small CRUD application that manages books

## Running

1. Copy .env.example as .env, optionally change values to your taste
2. Run `docker-compose up -d`
3. Hit `localhost:8085/book` endpoint with POST, GET, PUT and DELETE methods

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