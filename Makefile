update-deps:
	rm -rf vendor
	dep ensure
	dep ensure -update

test:
	go test -v