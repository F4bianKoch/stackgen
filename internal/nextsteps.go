package internal

import "fmt"

func PrintNextSteps(projectPath string) {
	fmt.Printf("\nNext Steps: \n")
	fmt.Printf(" - cd %s\n", projectPath)
	fmt.Printf(" - docker compose up -d\n")
	fmt.Printf(" - docker compose logs -f\n")
}
