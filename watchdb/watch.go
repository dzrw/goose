package watchdb

type Watch struct {
	DataSourceName string
	Method         string
	Path           string
	Echo           Echo

	Tag string
	id  int
}

func NewWatch(dsn, path, method, tag string, echo Echo) *Watch {
	return &Watch{
		DataSourceName: dsn,
		Method:         method,
		Path:           path,
		Tag:            tag,
		Echo:           echo,
		id:             -1,
	}
}
