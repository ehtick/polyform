package manifest

import "io"

type Artifact interface {
	Write(io.Writer) error
	Mime() string
}
