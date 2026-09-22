package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ServiceProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewServiceProxy(targetURL string) (*ServiceProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	return &ServiceProxy{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}, nil
}

func (p *ServiceProxy) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	p.proxy.ServeHTTP(w, r)
}
