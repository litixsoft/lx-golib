package lxRepo

import (
	"github.com/litixsoft/lx-golib/db"
)

type IBaseRepo interface {
	List(result interface{}, opts *lxDb.Options) (int, error)
}

type BaseRepo struct {
	baseDb lxDb.IBaseDb
}

func NewBaseRepo(db lxDb.IBaseDb) BaseRepo {
	return BaseRepo{baseDb:db}
}

func (repo *BaseRepo) List(result interface{}, opts *lxDb.Options) (int, error) {
	n, err := repo.baseDb.GetAll(nil, result, opts)
	if err != nil {
		return n, err
	}

	return n, nil
}
