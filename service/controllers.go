package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-faster/jx"
	"github.com/maqdev/go-be-template/util/idutil"
	"github.com/maqdev/go-be-template/util/logutil"

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

func (h handler) AuthorsAuthorIDGet(ctx context.Context, params api.AuthorsAuthorIDGetParams) (*api.Author, error) {
	if params.AuthorID != 1 {
		return nil, ErrNotFound
	}
	return &api.Author{
		ID:   1,
		Name: "Mag",
		Extra: api.NewOptAuthorExtra(map[string]jx.Raw{
			"foo": jx.Raw(`{"bar": "baz"}`),
		}),
	}, nil
}

func (h handler) AuthorsGet(ctx context.Context, params api.AuthorsGetParams) (*api.PagedAuthors, error) {
	return nil, errors.New("abc")
}

func (h handler) AuthorsPost(ctx context.Context, req api.OptAuthorsPostReq) (*api.AuthorsPostCreated, error) {
	// TODO implement me
	panic("implement me")
}

func (h handler) LoginPost(ctx context.Context, req *api.LoginPostReq) (*api.LoginPostCreated, error) {
	// TODO implement me
	panic("implement me")
}

func (h handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	if errors.Is(err, ErrInvalidToken) {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response: api.Error{
				Code:    "INVALID_TOKEN",
				Message: "Invalid token",
			},
		}
	}

	if errors.Is(err, ErrNotFound) {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response: api.Error{
				Code:    "NOT_FOUND",
				Message: "Not found",
			},
		}
	}

	errorID := idutil.NewID()
	logutil.Get(ctx).Error("Unhandled error", "err", err, "id", errorID)

	return &api.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.Error{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Error #" + errorID,
		},
	}
}
