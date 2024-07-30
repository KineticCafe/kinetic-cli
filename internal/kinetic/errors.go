package kinetic

import (
	"fmt"
)

// An ExitCodeError indicates the main program should exit with the given
// code.
type ExitCodeError int

func (e ExitCodeError) Error() string {
	return fmt.Sprintf("exit status %d", int(e))
}
