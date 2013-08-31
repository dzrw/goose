package eventdb

type nopdb struct{}

func (*nopdb) Dial() error {
	return nil
}

func (*nopdb) Close() {
	return
}
