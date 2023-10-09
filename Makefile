lint:
	go fmt . && go vet

test:
	go -v test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

bench:
	go test -bench=Bench
# -count 5
# -benchtime=10s

dep:
	go mod download

vet:
	go vet

.PHONY: lint test test_coverage bench dep vet