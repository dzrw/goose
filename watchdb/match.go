package watchdb

import "errors"

var ErrNoMatch = errors.New("IsMatch == false")

var EmptyMatchData = &match{0, "unmatched", nil}

type match struct {
	id   int
	tag  string
	echo Echo
}

func NewMatchData(w *Watch) MatchData {
	return &match{w.id, w.Tag, w.Echo}
}

func (m *match) Id() int {
	return m.id
}

func (m *match) Tag() string {
	m.assertIsMatch()
	return m.tag
}

func (m *match) Echo() Echo {
	m.assertIsMatch()
	return m.echo
}

func (m *match) IsMatch() bool {
	return m.id != 0
}

func (m *match) assertIsMatch() {
	if !m.IsMatch() {
		panic(ErrNoMatch)
	}
}
