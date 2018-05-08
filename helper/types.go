package lxHelper

type M map[string]interface{}

//type PostRequest struct {
//	Data interface{}
//	Options lxDb.Options `json:"opts"`
//	Query   interface{}   `json:"query"`
//}
//
//func NewRequestOpts(opts string) (*RequestOpts, error) {
//	var data RequestOpts
//
//	if len(opts) > 0 {
//		err := json.Unmarshal([]byte(opts), &data)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &data, nil
//}
