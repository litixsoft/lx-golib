package lxRepo

import (
	"github.com/litixsoft/lx-golib/db"
)

type IBaseRepo interface {
	List(query, result interface{}, opts *lxDb.Options) (int, error)
}

type BaseRepo struct {
	Db lxDb.IBaseDb
}

func NewBaseRepo(db lxDb.IBaseDb) *BaseRepo {
	return &BaseRepo{Db:db}
}

func (repo *BaseRepo) List(query, result interface{}, opts *lxDb.Options) (int, error) {
	n, err := repo.Db.GetAll(query, result, opts)
	if err != nil {
		return n, err
	}

	return n, nil
}
