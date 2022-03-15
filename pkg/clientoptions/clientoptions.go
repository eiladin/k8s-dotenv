package clientoptions

type Clientoptions struct {
	Namespace    string
	ShouldExport bool
}

func New() *Clientoptions {
	return &Clientoptions{}
}
