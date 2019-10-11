## install migrate cli
brew install golang-migrate
migrate -source file://path/to/migrations -database postgres://localhost:5432/database up
## docker
docker run -v /Users/apple/Documents/otus/calendar/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://calendar:123456@localhost:5432/calendar?sslmode=disable up
## docker-compose
docker-compose run migrations -path /migrations/ -database postgres://calendar:123456@localhost:5432/calendar?sslmode=disable up