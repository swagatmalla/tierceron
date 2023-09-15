//go:build darwin
// +build darwin

package memprotectopts

import (
	"log"

	"github.com/trimble-oss/tierceron/utils/mlock"
	"golang.org/x/sys/unix"
)

func MemProtectInit(logger *log.Logger) error {
	mlock.Mlock(logger)
}

func MemUnprotectAll(logger *log.Logger) error {
	return unix.Munlockall()
}

func MemProtect(logger *log.Logger, sensitive *string) error {
	return mlock.Mlock(logger*log.Logger, sensitive*string)
}