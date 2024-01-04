package service

import (
	"context"
	"github.com/maqdev/go-be-template/config"
	api "github.com/maqdev/go-be-template/gen/api/authors"
)

func NewHandler(cfg *config.AppConfig) api.Handler {
	return &handler{
		cfg: cfg,
	}
}

type handler struct {
	cfg *config.AppConfig
}

func (h handler) AuthorsAuthorIdGet(ctx context.Context, params api.AuthorsAuthorIdGetParams) (*api.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (h handler) AuthorsGet(ctx context.Context, params api.AuthorsGetParams) (*api.PagedAuthors, error) {
	//TODO implement me
	panic("implement me")
}

func (h handler) AuthorsPost(ctx context.Context, req api.OptAuthorsPostReq) (*api.AuthorsPostCreated, error) {
	//TODO implement me
	panic("implement me")
}

func (h handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	//TODO implement me
	panic("implement me")
}
