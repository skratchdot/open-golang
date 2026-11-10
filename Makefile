all: test check

.PHONY: test
test:
	go test -v --count=1 ./...

.PHONY: check
check:
	nix flake check

readme:
	godoc -templates=./readme-template/ "github.com/skratchdot/open-golang/open" > ./README.md