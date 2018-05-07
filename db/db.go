package lxDb

// ChangeInfo holds details about the outcome of an update operation.
type ChangeInfo struct {
	Updated int // Number of documents updated
	Removed int // Number of documents removed
	Matched int // Number of documents matched but not necessarily changed
}

type Options struct {
	Skip int `json:"skip"`
	Limit int `json:"limit"`
	Count bool `json:"count"`
}

// Base interface
type IBaseDb interface {
	Setup(config interface{}) error
	Create(data interface{}) error
	GetAll(query interface{}, result interface{}, opts *Options) (int, error)
	GetCount(query interface{}) (int, error)
	GetOne(query interface{}, result interface{}) error
	Update(query interface{}, data interface{}) error
	UpdateAll(query interface{}, data interface{}) (ChangeInfo, error)
	Delete(query interface{}) error
	DeleteAll(query interface{}) (ChangeInfo, error)
}