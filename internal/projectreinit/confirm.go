package internal

import (
	"fmt"
	"strings"
)

func confirm() bool {
	var permission string

	fmt.Printf("This will reset the current StackGen project.\n")
	fmt.Printf("All generated files may be overwritten.\n")
	fmt.Printf("Are you sure you want to continue? [y/N]: ")
	fmt.Scanln(&permission)
	fmt.Printf("\n")

	permission = strings.ToLower(permission)

	switch permission {
	case "y", "yes":
		return true
	default:
		fmt.Println("Reinitialization cancelled!")
		return false
	}
}
