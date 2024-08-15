// Пакет reqtype содержит тип для work-poller
package reqtype

import "github.com/go-resty/resty/v2"

// Тип ReqType содержит в себе endpoint и request
type ReqType struct {
	Req      *resty.Request
	Endpoint string
}
