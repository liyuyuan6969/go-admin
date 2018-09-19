package http

import (
	"bytes"
	"errors"
	"github.com/chenhg5/go-admin/context"
	"github.com/chenhg5/go-admin/plugins"
	"net/http"
	"strings"
	"github.com/chenhg5/go-admin/engine"
)

type Http struct {
}

func init()  {
	engine.Register(new(Http))
}

func (gins *Http) Use(router interface{}, plugin []plugins.Plugin) error {
	var (
		eng *http.ServeMux
		ok     bool
	)
	if eng, ok = router.(*http.ServeMux); !ok {
		return errors.New("wrong parameter")
	}

	var reqs map[string][]context.Path
	for _, plug := range plugin {
		reqs = ConstructNetHttpRequest(plug.GetRequest())
		for basicUrl, reqlist := range reqs {
			eng.HandleFunc(basicUrl, func(httpWriter http.ResponseWriter, httpRequest *http.Request) {
				for _, req := range reqlist {
					if httpRequest.Method == strings.ToUpper(req.Method) {
						ctx := context.NewContext(httpRequest)
						plug.GetHandler(req.URL, req.Method)(context.NewContext(httpRequest))
						httpWriter.WriteHeader(ctx.Response.StatusCode)
						if ctx.Response.Body != nil {
							buf := new(bytes.Buffer)
							buf.ReadFrom(ctx.Response.Body)
							httpWriter.Write(buf.Bytes())
						}
					}
				}
			})
		}
	}

	return nil
}

func ConstructNetHttpRequest(reqs []context.Path) map[string][]context.Path {
	var (
		NetHttpRequest = make(map[string][]context.Path, 0)
		usedUrl        []string
	)
	for _, req := range reqs {
		for _, url := range usedUrl {
			if url == req.URL {
				continue
			}
		}
		usedUrl = append(usedUrl, req.URL)
		NetHttpRequest[req.URL] = append(NetHttpRequest[req.URL], req)
	}
	return NetHttpRequest
}