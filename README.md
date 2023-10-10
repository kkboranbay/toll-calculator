# toll-calculator
A toll calculation service typically refers to a service or system that helps individuals or businesses calculate the cost of using toll roads, bridges, tunnels, or other tolled infrastructure. Toll roads and bridges often require drivers to pay a fee or toll to use them, and the amount of the toll can vary depending on factors such as the distance traveled, the type of vehicle, and any applicable discounts or promotions.

In the context of toll collection and road usage fees, an On-Board Unit (OBU) is a device installed in a vehicle.

## Installing protobuf compiler (protoc compiler)
For linux users or (WSL2)
```
sudo apt install -y protobuf-compiler
```

For Mac users you can use Brew for this
```
brew install protobuff
```

## Installing gRPC and Protobuffer plugins for Golang.
1. Protobuffers
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

2. gRPC
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

3. Note that you need to set the /go/bin directory in your path
```
export PATH="$PATH:$(go env GOPATH)/bin"
```

4. Install the package dependencies
```
go get google.golang.org/protobuf
```
```
go get google.golang.org/grpc
```