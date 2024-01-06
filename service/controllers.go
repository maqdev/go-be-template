package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maqdev/go-be-template/gen/db"

	"github.com/go-faster/jx"
	"github.com/maqdev/go-be-template/util/idutil"
	"github.com/maqdev/go-be-template/util/logutil"

	"github.com/maqdev/go-be-template/config"
	api "github.com/maqdev/go-be-template/gen/api/authors"
)

const MaxPageSize = 500
const DefaultPageSize = 10

func NewHandler(cfg *config.AppConfig, dbPool *pgxpool.Pool) api.Handler {
	return &handler{
		cfg:     cfg,
		queries: db.New(dbPool),
	}
}

type handler struct {
	cfg     *config.AppConfig
	queries *db.Queries
}

func (h *handler) AuthorsAuthorIDGet(ctx context.Context, params api.AuthorsAuthorIDGetParams) (*api.Author, error) {
	author, err := h.queries.GetAuthor(ctx, params.AuthorID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get author from DB: %w", err)
	}

	authorRes, err := dbAuthorToResp(author)
	if err != nil {
		return nil, fmt.Errorf("failed to convert author to response: %w", err)
	}

	return authorRes, nil
}

func (h handler) AuthorsGet(ctx context.Context, params api.AuthorsGetParams) (*api.PagedAuthors, error) {
	listParams := db.ListAuthorsParams{
		Previd: pgtype.Int8{},
		Lim:    min(params.Limit.Or(DefaultPageSize), MaxPageSize),
	}
	if params.Token.Set {
		prevID, err := strconv.ParseInt(params.Token.Value, 36, 64)
		if err == nil {
			listParams.Previd.Int64 = prevID
			listParams.Previd.Valid = true
		}
	}

	authors, err := h.queries.ListAuthors(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list authors from DB: %w", err)
	}

	res := &api.PagedAuthors{
		Content:   nil,
		NextToken: api.OptString{},
	}
	if len(authors) > 0 {
		res.Content = make([]api.Author, len(authors))
		for i, author := range authors {
			var resAuthor *api.Author
			resAuthor, err = dbAuthorToResp(author)
			if err != nil {
				return nil, fmt.Errorf("failed to convert author to response: %w", err)
			}
			res.Content[i] = *resAuthor
		}
		if len(authors) >= int(listParams.Lim) {
			res.NextToken.Value = strconv.FormatInt(authors[len(authors)-1].ID, 36)
			res.NextToken.Set = true
		}
	}
	return res, nil
}

func (h handler) AuthorsPost(ctx context.Context, req *api.AuthorsPostReq) (*api.AuthorsPostCreated, error) {
	// todo: validate name

	id, err := h.queries.CreateAuthor(ctx, db.CreateAuthorParams{
		Name: req.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create author in DB: %w", err)
	}

	return &api.AuthorsPostCreated{
		ID: id,
	}, nil
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

func dbAuthorToResp(author db.Author) (*api.Author, error) {
	var authorExtra api.OptAuthorExtra
	if author.Extra != nil {
		extra := make(map[string]jx.Raw)
		d := jx.DecodeBytes(author.Extra)
		err := d.Obj(func(d *jx.Decoder, key string) error {
			var err error
			extra[key], err = d.Raw()
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("failed to decode extra: %w", err)
		}
		authorExtra = api.NewOptAuthorExtra(extra)
	}

	return &api.Author{
		ID:    author.ID,
		Name:  author.Name,
		Extra: authorExtra,
	}, nil
}
