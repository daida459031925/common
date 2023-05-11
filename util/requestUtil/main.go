package requestUtil

import "net/http"

func GetRequest(reqs ...*http.Request) *http.Request {
	var req *http.Request = nil
	if reqs != nil && len(reqs) > 0 {
		req = reqs[0]
	}
	return req
}
