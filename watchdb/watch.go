package watchdb

type Watch struct {
	Method string
	Path   string
	Echo   Echo

	Tag string
	id  int
}

func NewWatch(path, method, tag string, echo Echo) *Watch {
	return &Watch{
		Method: method,
		Path:   path,
		Tag:    tag,
		Echo:   echo,
		id:     -1,
	}
}
