package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"web_app/controllers/rpc/proto"
	"web_app/logic"
)

type GetPostDetailServiceImpl struct{}

func (c *GetPostDetailServiceImpl) GetPostDetail(ctx context.Context, in *proto.PostID) (*proto.PostDetail, error) {
	post, err := logic.GetPostDetail(in.PostId)
	post_detail := &proto.PostDetail{
		AuthorName: post.AuthorName,
		VoteNum:    post.VoteNum,
		AuthorID:   post.AuthorID,
		Status:     post.Status,
		Title:      post.Title,
		Content:    post.Content,
	}

	return post_detail, err
}

func RegisterAndServe() {
	grpcServer := grpc.NewServer()
	proto.RegisterPostServiceServer(grpcServer, &GetPostDetailServiceImpl{})

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
	return
}

//
//func (c *GetPostDetailServiceImpl) GetPostDetail(ctx context.Context, in *proto.PostID, opts ...grpc.CallOption) (*proto.Post, error) {
//	post, err := logic.GetPostDetail(in.Id)
//	post_pro := &proto.Post{
//		ID:         post.ID,
//		AuthorID:   post.AuthorID,
//		Status:     post.Status,
//		Title:      post.Title,
//		Content:    post.Content,
//		CreateTime: post.CreateTime,
//	}
//
//	return post_pro, err
//}
