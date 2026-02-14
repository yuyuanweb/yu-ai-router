package yuairoutersdk

type FuncStreamCallback struct {
	OnMessageFunc  func(chunk ChatChunk)
	OnCompleteFunc func()
	OnErrorFunc    func(err error)
}

func (f FuncStreamCallback) OnMessage(chunk ChatChunk) {
	if f.OnMessageFunc != nil {
		f.OnMessageFunc(chunk)
	}
}

func (f FuncStreamCallback) OnComplete() {
	if f.OnCompleteFunc != nil {
		f.OnCompleteFunc()
	}
}

func (f FuncStreamCallback) OnError(err error) {
	if f.OnErrorFunc != nil {
		f.OnErrorFunc(err)
	}
}
