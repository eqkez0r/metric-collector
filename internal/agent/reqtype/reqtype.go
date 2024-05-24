package reqtype

import "github.com/go-resty/resty/v2"

type ReqType struct {
	Req      *resty.Request
	Endpoint string
}
