package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/go-kratos/kratos/v2/log"

	pb "blog/api/blog/v1/article"

	"blog/internal/biz"
)

type ArticleService struct {
	pb.UnimplementedArticleServiceServer

	log *log.Helper

	article *biz.ArticleUsecase
}

func NewArticleService(article *biz.ArticleUsecase, logger log.Logger) *ArticleService {
	return &ArticleService{
		article: article,
		log:     log.NewHelper(log.With(logger, "model", "service/article")),
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleReply, error) {
	s.log.Infof("input data %v", req)

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	err = s.article.Create(ctx, &biz.Article{
		Title:   req.Title,
		Content: req.Content,
	})
	return &pb.CreateArticleReply{}, err
}

func (s *ArticleService) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleReply, error) {
	s.log.Infof("input data %v", req)

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	err = s.article.Update(ctx, req.Id, &biz.Article{
		Title:   req.Title,
		Content: req.Content,
	})
	return &pb.UpdateArticleReply{}, err
}

func (s *ArticleService) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleReply, error) {
	s.log.Infof("input data %v", req)

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	err = s.article.Delete(ctx, req.Id)
	return &pb.DeleteArticleReply{}, err
}

func (s *ArticleService) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleReply, error) {
	tr := otel.Tracer("api")
	_, span := tr.Start(ctx, "GetArticle")
	defer span.End()
	p, err := s.article.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetArticleReply{Article: &pb.Article{Id: p.Id, Title: p.Title, Content: p.Content, Like: p.Like}}, nil
}

func (s *ArticleService) ListArticle(ctx context.Context, req *pb.ListArticleRequest) (*pb.ListArticleReply, error) {
	ps, err := s.article.List(ctx, req)
	reply := &pb.ListArticleReply{}
	for _, p := range ps {
		reply.Results = append(reply.Results, &pb.Article{
			Id:      p.Id,
			Title:   p.Title,
			Content: p.Content,
		})
	}
	return reply, err
}
