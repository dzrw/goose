package watchdb

type MatchExpr struct {
	Method string
	Path   string
}

func NewMatchExpr(path, method string) *MatchExpr {
	return &MatchExpr{Path: path, Method: method}
}
