package formatter

import (
	"fmt"

	"github.com/MikelMelnichuk/mycal/internal/models"
)

func PrintJSON(events []models.Event) {
	fmt.Println("Json Print")
	for _, e := range events {
		fmt.Printf("Title: %s, time %s-%s", e.Title, e.Start, e.End)
	}
}
