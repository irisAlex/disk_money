package mongodb

import "errors"

const (
	NoResultErr = "mongo: no documents in result"
	NotFound    = "Sorry, not found result in mongodb"
)

// mustHave ： 是否必须有搜索到数据
func HandleMongoErr(mustHave bool, err error, length int) error {
	if mustHave {
		if length == 0 {
			return errors.New(NotFound)
		}
		if err != nil && err.Error() == NoResultErr {
			return errors.New(NotFound)
		}
	} else {
		if err != nil && err.Error() == NoResultErr {
			return nil
		}
	}

	return err
}
