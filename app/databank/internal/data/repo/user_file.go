package repo

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"store/app/databank/internal/data/repo/ent"
	"store/pkg/clients"
)

type UserFileRepo struct {
	db    *ent.Client
	redis *clients.RedisClient
}

func NewUserFileRepo(db *ent.Client, redis *clients.RedisClient) *UserFileRepo {
	return &UserFileRepo{db: db, redis: redis}
}

func (t *UserFileRepo) InsertBulk(ctx context.Context, files []*ent.UserFile) error {

	if len(files) == 0 {
		return nil
	}

	var creates []*ent.UserFileCreate
	for _, x := range files {
		creates = append(creates,
			t.db.UserFile.Create().
				SetUserID(x.UserID).
				SetName(x.Name).
				SetMd5(x.Md5).
				SetName(x.Name).
				SetCategory(x.Category).
				SetSize(x.Size),
		)
	}

	err := t.db.UserFile.CreateBulk(creates...).OnConflict(sql.ResolveWithNewValues()).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
