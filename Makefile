build:
	@go build -o bin/gobank

run : build
	@./bin/gobank

test : 
	@go test -v ./...

# @ will not print out what steps are going. You can also remove @ and try
# 