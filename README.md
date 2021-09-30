# Booker

A small CRUD application that manages books

## Testing
Run `make test`

## Running

### Docker-compose
1. Copy .env.example as .env, optionally change values to your taste
2. Run `make run-docker`
3. Hit `localhost:8080/book` endpoint with POST, GET, PUT and DELETE methods. Change 8080 to the value of APP_PORT in .env file, if you changed it.

### Locally
1. Run `make build`
2. Run `make run-local`


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