package output

import (
	"context"
	"io"
)

type Actor interface {
	Exec(ctx context.Context, stdout, stderr io.Writer) error
}

type Config interface {
	RequireEmbededThisMethod()
}

type configBase struct{}

func (*configBase) RequireEmbededThisMethod() {}
