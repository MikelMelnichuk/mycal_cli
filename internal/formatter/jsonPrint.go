package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/MikelMelnichuk/mycal/internal/models"
)

func PrintJSON(events []models.Event) {
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		fmt.Printf("Could not convert to json %s\n", err)
	}
	fmt.Println(string(data))
}
