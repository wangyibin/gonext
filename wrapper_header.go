package gonext

import (
	echoEngine "github.com/labstack/echo/engine"
)

type (
	// Header implements `engine.Header`.
	headerWrapperFromEcho struct {
		echoEngine.Header
	}
)

// Add implements `engine.Header#Add` function.
func (h *headerWrapperFromEcho) Add(key, val string) {
	h.Header.Add(key, val)
}

// Del implements `engine.Header#Del` function.
func (h *headerWrapperFromEcho) Del(key string) {
	h.Header.Del(key)
}

// Set implements `engine.Header#Set` function.
func (h *headerWrapperFromEcho) Set(key, val string) {
	h.Header.Set(key, val)
}

// Get implements `engine.Header#Get` function.
func (h *headerWrapperFromEcho) Get(key string) string {
	return h.Header.Get(key)
}

// Keys implements `engine.Header#Keys` function.
func (h *headerWrapperFromEcho) Keys() []string {
	return h.Header.Keys()
}
