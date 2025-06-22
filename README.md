# Cache with redis

This project demonstrates a basic caching mechanism implemented in Go, leveraging Redis as the cache store and SQLite for data persistence. It's designed to provide a straightforward introduction to caching principles and their practical application.

## Requirements:
- go 1.24.4+
- docker
- docker compose


## How to run:

1. Clone the repository
```bash
git clone https://github.com/LXSCA7/redis-golang
cd redis-golang
```

2. Get dependencies
```bash
go mod tidy
```

3. Start the Redis Container
```bash
docker compose up --build -d
```

**To verify Redis is running and responsive, you can use redis-cli:** `docker exec -it redis-db redis-cli PING`

4. Run the code
First time (for populate the database):
```bash
go run . insert
```

Subsequent runs:
```bash
go run .