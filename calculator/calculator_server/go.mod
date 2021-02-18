module example.com/calculator_server

go 1.16

replace example.com/calculatorpb => ../calculatorpb

require (
	example.com/calculatorpb v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.35.0
)
