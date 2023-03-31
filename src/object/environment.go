package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	curEnv := e
	for {
		_, ok := curEnv.store[name]
		if ok {
			break
		}
		if curEnv.outer == nil {
			curEnv = e
			break
		}
		curEnv = curEnv.outer
	}
	curEnv.store[name] = val
	return val
}

func (e *Environment) SetCurrentEnv(name string, val Object) Object {
	e.store[name] = val
	return val
}