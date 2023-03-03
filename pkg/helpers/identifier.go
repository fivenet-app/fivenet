package helpers

import "fmt"

func BuildCharSearchIdentifier(license string) string {
	return fmt.Sprintf("char%%:%s", license)
}
