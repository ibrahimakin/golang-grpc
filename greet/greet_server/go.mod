module example.com/greet_server

go 1.15

replace example.com/greetpb => ../greetpb

require (
	example.com/greetpb v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.35.0
)
