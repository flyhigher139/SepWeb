package context

import "strconv"

type StringValue struct {
	val string
	err error
}

func (s StringValue) ToString() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.val, nil
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
