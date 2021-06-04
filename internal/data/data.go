package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"

	"blog/internal/conf"
	"blog/internal/data/ent"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewArticleRepo, NewCommentRepo)

// Data .
type Data struct {
	db  *ent.Client
	rdb *redis.Client
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger) (*Data, error) {
	helper := log.NewHelper(log.With(logger, "model", "data/article"))

	// create database
	defaultSource := fmt.Sprintf(conf.Database.Source, conf.Database.Driver)
	db, err := sql.Open(conf.Database.Driver, defaultSource)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	execArgs := fmt.Sprintf("CREATE DATABASE %v", conf.Database.Dbname)
	_, err = db.Exec(execArgs)
	if err != nil {
		// database exists
		helper.Infof("create database [%s] err: %s", conf.Database.Dbname, err.Error())
	} else {
		helper.Infof("create database [%s] success", conf.Database.Dbname)
	}

	// create db client
	source := fmt.Sprintf(conf.Database.Source, conf.Database.Dbname)
	client, err := ent.Open(
		conf.Database.Driver,
		source,
	)
	if err != nil {
		helper.Errorf("failed opening connection to database: %v", err)
		return nil, err
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		helper.Errorf("failed creating schema resources: %v", err)
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	return &Data{
		db:  client,
		rdb: rdb,
	}, nil
}
