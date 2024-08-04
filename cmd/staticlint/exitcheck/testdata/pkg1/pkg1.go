package main

import "os"

func Exit(_ int) {
}

func main() {
	Exit(0)
	os.Exit(1) // want "calling os.Exit func in main func"
}
