//go:build windows

package store

func setUmask(int) int { return 0 }
