package repo

import "store/pkg/clients"

type UserRepo struct {
	db *clients.ClickHouseClient
}

func NewUserRepo(db *clients.ClickHouseClient) *UserRepo {
	return &UserRepo{db: db}
}

func (t *UserRepo) FindByEmail() {

}
