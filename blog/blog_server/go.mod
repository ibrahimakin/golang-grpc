module example.com/blog_server

go 1.16

replace example.com/blogpb => ../blogpb

require (
	example.com/blogpb v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.4.6
	google.golang.org/grpc v1.35.0
)
