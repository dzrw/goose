package eventdb

type fakeEventProvider struct{}

func (*fakeEventProvider) Dial() error {
	return nil
}

func (*fakeEventProvider) Close() {}
