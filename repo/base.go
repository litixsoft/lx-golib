package lxRepo

import (
	"github.com/litixsoft/lx-golib/db"
)

type IBaseRepo interface {
	List(result interface{}, opts *lxDb.Options) (int, error)
}

type BaseRepo struct {
	Db lxDb.IDatabase
}

func (repo *BaseRepo) List(result interface{}, opts *lxDb.Options) (int, error) {
	n, err := repo.Db.GetAll(nil, result, opts)
	if err != nil {
		return n, err
	}

	return n, nil
}
