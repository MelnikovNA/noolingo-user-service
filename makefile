refact:
	go mod edit -module github.com/noolingo/user-service
	-- rename all imported module
	find . -type f -name '*.go' \
  	-exec sed -i -e 's,github.com/MelnikovNA/noolingo-user-service,github.com/noolingo/user-service,g' {} \;