package watch

import "errors"

// ErrNoCommand is returned when Watch is called with an empty args slice.
var ErrNoCommand = errors.New("watch: no command specified")
