gate:
	@go build -o bin/gate gateway/main.go
	@./bin/gate

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data-receiver
	@./bin/receiver	

calculator:
	@go build -o bin/calculator ./distance-calculator
	@./bin/calculator	

agg:
	@go build -o bin/agg ./aggregator
	@./bin/agg	

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

# if you get error -> make: `obu' is up to date.
.PHONY: obu
# Using .PHONY is a way to avoid conflicts and confusion 
# that might arise if there happened to be a file named "obu" 
# in the same directory, which could potentially interfere with the target. 
# By declaring it as a phony target, you make it clear that it's a special 
# target meant to execute specific commands.