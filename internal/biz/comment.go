package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	pb "blog/api/blog/v1/comment"
)

type Comment struct {
	Id        int64
	Name      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Like      int64
}

type CommentRepo interface {
	// db
	ListComment(ctx context.Context, req *pb.ListCommentRequest) ([]*Comment, error)
	GetComment(ctx context.Context, id int64) (*Comment, error)
	CreateComment(ctx context.Context, comment *Comment) error
	UpdateComment(ctx context.Context, id int64, comment *Comment) error
	DeleteComment(ctx context.Context, id int64) error

	// redis
	GetCommentLike(ctx context.Context, id int64) (rv int64, err error)
	IncCommentLike(ctx context.Context, id int64) error
}

type CommentUsecase struct {
	repo CommentRepo
}

func NewCommentUsecase(repo CommentRepo, logger log.Logger) *CommentUsecase {
	return &CommentUsecase{repo: repo}
}

func (uc *CommentUsecase) List(ctx context.Context, req *pb.ListCommentRequest) (ps []*Comment, err error) {
	ps, err = uc.repo.ListComment(ctx, req)
	if err != nil {
		return
	}
	return
}

func (uc *CommentUsecase) Get(ctx context.Context, id int64) (p *Comment, err error) {
	p, err = uc.repo.GetComment(ctx, id)
	if err != nil {
		return
	}
	err = uc.repo.IncCommentLike(ctx, id)
	if err != nil {
		return
	}
	p.Like, err = uc.repo.GetCommentLike(ctx, id)
	if err != nil {
		return
	}
	return
}

func (uc *CommentUsecase) Create(ctx context.Context, comment *Comment) error {
	return uc.repo.CreateComment(ctx, comment)
}

func (uc *CommentUsecase) Update(ctx context.Context, id int64, comment *Comment) error {
	return uc.repo.UpdateComment(ctx, id, comment)
}

func (uc *CommentUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteComment(ctx, id)
}
