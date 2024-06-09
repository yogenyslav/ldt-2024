package pkg

import (
	"fmt"
)

// GetCallbackData returns a string with a special character and the text.
func GetCallbackData(text string) string {
	return fmt.Sprintf("\f%s", text)
}
