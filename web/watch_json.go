package web

import (
	"errors"
	"github.com/drone/routes"
	"github.com/politician/goose/watchdb"
	"net/http"
)

type JsonWatch struct {
	Tag            string
	DataSourceName string
	MatchExpr      *JsonMatchExpr
	Echo           *JsonEcho
}

type JsonMatchExpr struct {
	Method string
	Path   string
}

type JsonEcho struct {
	Status  int
	Headers map[string]string
	Body    string
}

func ParseJsonWatch(req *http.Request) (w *watchdb.Watch, err error) {
	jw := &JsonWatch{}
	err = routes.ReadJson(req, jw)
	if err != nil {
		return
	}

	if jw.DataSourceName == "" {
		err = errors.New("missing provider")
		return
	}

	path := jw.MatchExpr.Path
	method := jw.MatchExpr.Method
	tag := jw.Tag
	echo := ParseJsonEcho(jw.Echo)

	w = watchdb.NewWatch(path, method, tag, echo)
	return
}

func ParseJsonEcho(j *JsonEcho) watchdb.Echo {
	return &echo{
		status:  j.Status,
		headers: j.Headers,
		body:    []byte(j.Body),
	}
}
