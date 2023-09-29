package extras

type Memory interface {
	Get(idx interface{}) interface{}
	Set(idx interface{}, val interface{})
	LoadProgram(programName string)
}
