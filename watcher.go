package main

import "time"

// WatchFunc is the signature for a trigger watcher. It blocks until exit is
// closed or an unrecoverable error occurs, calling onOpen/onClose on each
// relevant state transition.
type WatchFunc func(exit <-chan struct{}, onOpen, onClose func()) error

// pollInterval is how often watchers poll their underlying sensor or device.
const pollInterval = 2 * time.Second
