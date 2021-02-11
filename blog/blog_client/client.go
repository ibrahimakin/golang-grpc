package main

import (
	"context"
	"fmt"
	"log"

	"../blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Blog client")

	opts := grpc.WithInsecure()
	tls := false
	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust Certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v\n", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close() // Finally close the connection

	c := blogpb.NewBlogServiceClient(cc)

	// Create Blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "ibrahim",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// Read Blog
	fmt.Println("Reading the blog")
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "12edadw3dad"})
	if err2 != nil {
		log.Printf("Error happend when reading: %v\n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)

	if readBlogErr != nil {
		log.Fatalf("Error happened when reading: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// Update Blog
	fmt.Println("Updating the blog")
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (Edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateBlogRes, updateBlogErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateBlogErr != nil {
		log.Fatalf("Unexpected error: %v\n", updateBlogErr)
	}
	fmt.Printf("Blog has been updated: %v\n", updateBlogRes)

	// Delete Blog
	fmt.Println("Deleting the blog")
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		log.Fatalf("Error happened when deleting: %v\n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v\n", deleteRes)
}
