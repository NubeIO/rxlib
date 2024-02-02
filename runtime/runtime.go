package runtime

import "github.com/NubeIO/rxlib"

type Runtime interface {
	Get() map[string]rxlib.Object
}
