package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"

	pb "blog/api/blog/v1/comment"
	"blog/internal/biz"
)

type commentRepo struct {
	data *Data
	log  *log.Helper
}

// NewCommentRepo .
func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "model", "data/comment")),
	}
}

func (cr *commentRepo) ListComment(ctx context.Context, req *pb.ListCommentRequest) ([]*biz.Comment, error) {
	ps, err := cr.data.db.Comment.
		Query().
		Limit(int(req.Limit)).
		Offset(int(req.Offset)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.Comment, 0)
	for _, p := range ps {
		rv = append(rv, &biz.Comment{
			Id:        p.ID,
			Name:      p.Name,
			Content:   p.Content,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}
	return rv, nil
}

func (cr *commentRepo) GetComment(ctx context.Context, id int64) (*biz.Comment, error) {
	p, err := cr.data.db.Comment.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.Comment{
		Id:        p.ID,
		Name:      p.Name,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}

func (cr *commentRepo) CreateComment(ctx context.Context, comment *biz.Comment) error {
	_, err := cr.data.db.Comment.
		Create().
		SetName(comment.Name).
		SetContent(comment.Content).
		Save(ctx)
	return err
}

func (cr *commentRepo) UpdateComment(ctx context.Context, id int64, comment *biz.Comment) error {
	p, err := cr.data.db.Comment.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = p.Update().
		SetName(comment.Name).
		SetContent(comment.Content).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	return err
}

func (cr *commentRepo) DeleteComment(ctx context.Context, id int64) error {
	return cr.data.db.Comment.DeleteOneID(id).Exec(ctx)
}

// -----------------redis-------------------

func (cr *commentRepo) GetCommentLike(ctx context.Context, id int64) (rv int64, err error) {
	get := cr.data.rdb.Get(ctx, likeKey(id))
	rv, err = get.Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return
}

func (cr *commentRepo) IncCommentLike(ctx context.Context, id int64) error {
	_, err := cr.data.rdb.Incr(ctx, likeKey(id)).Result()
	return err
}
