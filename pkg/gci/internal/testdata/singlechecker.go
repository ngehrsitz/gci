package proc

import (
	"context" // in-line comment

	//nolint:depguard // A multi-line comment explaining why in
	// this one case it's OK to use os/exec even though depguard
	// is configured to force us to use dlib/exec instead.
	"os/exec"
)

func main() {
	_ = context.Canceled
	_ = exec.ErrNotFound
}
