package tests__test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maqdev/go-be-template/config"
	"github.com/maqdev/go-be-template/gen/db"
	"github.com/stretchr/testify/require"
)

func Test_Create_And_GetAuthor(t *testing.T) {
	cfg := config.TestConfig(t)
	ctx := context.Background()
	req := require.New(t)
	dbPool, err := cfg.DB.CreatePool(ctx)
	req.NoError(err)
	defer dbPool.Close()

	queries := db.New(dbPool)

	name := faker.Name()
	tm := time.Now()

	id, err := queries.CreateAuthor(ctx, db.CreateAuthorParams{
		Name:  name,
		Extra: nil,
	})
	req.NoError(err)
	req.NotZero(id)

	var a db.Author
	a, err = queries.GetAuthor(ctx, id)
	req.NoError(err)

	req.Equal(id, a.ID)
	req.Equal(name, a.Name)
	req.Nil(a.Extra)
	req.Greater(a.CreatedAt.Time.Unix(), tm.Unix()-1000)
}

func Test_Create_And_GetAuthors(t *testing.T) {
	cfg := config.TestConfig(t)
	ctx := context.Background()
	req := require.New(t)
	dbPool, err := cfg.DB.CreatePool(ctx)
	req.NoError(err)
	defer dbPool.Close()

	queries := db.New(dbPool)

	total := 5
	ids := make([]int64, 0, total)

	for i := 0; i < total; i++ {
		name := faker.Name()
		var id int64
		id, err = queries.CreateAuthor(ctx, db.CreateAuthorParams{
			Name:  name,
			Extra: nil,
		})
		req.NoError(err)
		req.NotZero(id)
		ids = append(ids, id)
	}

	var authors []db.Author
	authors, err = queries.ListAuthors(ctx, db.ListAuthorsParams{
		Previd: pgtype.Int8{Int64: ids[0] - 1, Valid: true},
		Lim:    int32(100 + total),
	})
	req.NoError(err)
	req.GreaterOrEqual(len(authors), total)
	for i := 0; i < total; i++ {
		found := slices.IndexFunc(authors, func(author db.Author) bool {
			return author.ID == ids[i]
		}) >= 0

		req.True(found, "Author %d not found", ids[i])
	}
}
