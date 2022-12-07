package logs

type containerLogManagerStub struct{}

func (*containerLogManagerStub) Start() {}

// NewStubContainerLogManager returns an empty ContainerLogManager which does nothing.
func NewStubContainerLogManager() ContainerLogManager {
	return &containerLogManagerStub{}
}
