package version

import "fmt"

const Version = "v0.1.1"

func PrintVersion() {
	fmt.Println("stackgen", Version)
}
