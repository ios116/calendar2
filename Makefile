up:
	docker-compose up -d --build

down:
	docker-compose down

restart: down up

test:
	docker-compose -f docker-compose.test.yaml up --build -d;\
	docker-compose -f docker-compose.test.yaml run integration_tests go test -v -count=1 ./... || test_status=$$?;\
	docker-compose -f docker-compose.test.yaml down; echo "status="$$test_status;exit $$test_status;

