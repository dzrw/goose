package watchdb

type MatchExpr struct {
	DataSourceName string
	Method         string
	Path           string
}

func NewMatchExpr(dataSourceName, path, method string) *MatchExpr {
	return &MatchExpr{dataSourceName, method, path}
}
