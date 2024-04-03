package metric

import "strconv"

type Counter int64

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}
