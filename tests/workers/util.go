package workers

type multiCloser struct {
	fns []func() error
}

func (mc *multiCloser) F() func() error {
	return func() error {
		var err error
		for i := range mc.fns {
			if err1 := mc.fns[len(mc.fns)-1-i](); err == nil {
				err = err1
			}
		}
		mc.fns = nil
		return err
	}
}

func (mc *multiCloser) append(f func() error) {
	mc.fns = append(mc.fns, f)
}
