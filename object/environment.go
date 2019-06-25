package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	curdir, _ := env.Get("__curdir")
	env.Set("__curdir", curdir)
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnvironmentWithDir(curdir string) *Environment {
	env := NewEnvironment()
	env.Set("__curdir", &String{Value: curdir})
	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Keys() []string {
	keys := make([]string, len(e.store))
	i := 0
	for key := range e.store {
		keys[i] = key
		i++
	}
	return keys
}
