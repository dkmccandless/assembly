package object

type Environment struct{ store map[string]Object }

func NewEnvironment() *Environment { return &Environment{store: make(map[string]Object)} }

func (e *Environment) Get(s string) (Object, bool) {
	obj, ok := e.store[s]
	return obj, ok
}

func (e *Environment) Set(s string, obj Object) { e.store[s] = obj }
