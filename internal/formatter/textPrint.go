package formatter

import (
	"fmt"

	"github.com/MikelMelnichuk/mycal/internal/models"
)

func PrintText(events []models.Event) {
	fmt.Println("TEXT print")
	for _, e := range events {
		fmt.Printf("Title: %s, time %s-%s", e.Title, e.Start, e.End)
	}
}
