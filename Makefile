obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data-receiver
	@./bin/receiver	

calculator:
	@go build -o bin/calculator ./distance-calculator
	@./bin/calculator	

# if you get error -> make: `obu' is up to date.
.PHONY: obu
# Using .PHONY is a way to avoid conflicts and confusion 
# that might arise if there happened to be a file named "obu" 
# in the same directory, which could potentially interfere with the target. 
# By declaring it as a phony target, you make it clear that it's a special 
# target meant to execute specific commands.