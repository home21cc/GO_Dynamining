package castings

type castInterfaceFunc func() CastingInterface

type CastingInterface interface {
	Init( config string) error
	CastingValue()
	Destroy()
	Flush()
}

