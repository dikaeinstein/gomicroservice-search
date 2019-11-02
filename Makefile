unit:
	go test -v -race $(shell go list ./... | grep -v /vendor/)

staticcheck:
	staticcheck -tests=false $(shell go list ./... | grep -v /vendor/)

safesql:
	safesql github.com/dikaeinstein/gomicroservice-search/cmd/

benchmark:
	go test -benchmem -benchtime=20s -bench=. github.com/dikaeinstein/gomicroservice-search/handler | tee bench.txt
	if [ -a old_bench.txt ]; then \
  	benchcmp old_bench.txt bench.txt; \
	fi;

	if [ $$? -eq 0 ]; then \
		mv bench.txt old_bench.txt; \
	fi;

build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./search cmd/search.go

build_docker:build_linux
	docker build -t dikaeinstein/gomicroservice-search .

start_stack:
	docker-compose up -d

run: start_stack
	export MYSQL_CONNECTION="root:password@tcp(${DOCKER_IP}:3307)/kittens"; \
	export DOGSTATSD=127.0.0.1:8125; \
	go run cmd/search.go
	docker-compose stop

RSA_PUBLIC_KEY := $(RSA_PUBLIC_KEY)
DATADOG_API_KEY := $(DATADOG_API_KEY)
run_docker: start_stack
	docker run --network gomicroservice-search_default \
	-p 8082:8082 \
	-e "MYSQL_CONNECTION=root:password@tcp(mysql:3306)/kittens" \
	-e "DOGSTATSD=localhost:8125" \
	-e "DD_SITE=datadoghq.eu" \
	-e "DD_API_KEY=$$DATADOG_API_KEY" \
	-e "RSA_PUBLIC_KEY=$$RSA_PUBLIC_KEY" \
	dikaeinstein/gomicroservice-search:latest
	docker-compose stop

integration: start_stack
	cd features && MYSQL_CONNECTION="root:password@tcp(${DOCKER_IP}:3307)/kittens"  DOGSTATSD=localhost:8125 godog ./
	docker-compose stop
	docker-compose rm -f

test: unit benchmark staticcheck safesql integration

circleintegration:
	docker build -t circletemp -f ./Dockerfile.integration .
	docker-compose f docker-compose.start_stack.yml up -d
	docker run --network gomicroservice-search_default -w /go/src/github.com/dikaeinstein/gomicroservice-search/features -e "MYSQL_CONNECTION=root:password@tcp(mysql:3306)/kittens" circletemp godog ./
	docker-compose stop
	docker-compose rm -f
