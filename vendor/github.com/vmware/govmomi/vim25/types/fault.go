package types

type HasFault interface {
	Fault() BaseMethodFault
}

func IsFileNotFound(err error) bool {
	if f, ok := err.(HasFault); ok {
		switch f.Fault().(type) {
		case *FileNotFound:
			return true
		}
	}

	return false
}
