package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/go-kratos/kratos/v2/log"

	pb "blog/api/blog/v1/comment"

	"blog/internal/biz"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer

	log *log.Helper

	comment *biz.CommentUsecase
}

func NewCommentService(comment *biz.CommentUsecase, logger log.Logger) *CommentService {
	return &CommentService{
		comment: comment,
		log:     log.NewHelper("comment", logger),
	}
}

func (s *CommentService) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentReply, error) {
	s.log.Infof("input data %v", req)
	err := s.comment.Create(ctx, &biz.Comment{
		Name:   req.Name,
		Content: req.Content,
	})
	return &pb.CreateCommentReply{}, err
}

func (s *CommentService) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentReply, error) {
	s.log.Infof("input data %v", req)
	err := s.comment.Update(ctx, req.Id, &biz.Comment{
		Content: req.Content,
	})
	return &pb.UpdateCommentReply{}, err
}

func (s *CommentService) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentReply, error) {
	s.log.Infof("input data %v", req)
	err := s.comment.Delete(ctx, req.Id)
	return &pb.DeleteCommentReply{}, err
}

func (s *CommentService) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentReply, error) {
	tr := otel.Tracer("api")
	_, span := tr.Start(ctx, "GetComment")
	defer span.End()
	p, err := s.comment.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetCommentReply{Comment: &pb.Comment{Id: p.Id, Name: p.Name, Content: p.Content, Like: p.Like}}, nil
}

func (s *CommentService) ListComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.ListCommentReply, error) {
	ps, err := s.comment.List(ctx, req)
	reply := &pb.ListCommentReply{}
	for _, p := range ps {
		reply.Results = append(reply.Results, &pb.Comment{
			Id:      p.Id,
			Name:   p.Name,
			Content: p.Content,
		})
	}
	return reply, err
}
