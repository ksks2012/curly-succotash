package main

import (
	"curly-succotash/backend/internal/model"

	"encoding/json"
	"fmt"
)

type cardResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Effect      string `json:"effect"`
}

var storyResponse = `{
  "story_background": "The realm of Eldoria, nestled between the Whispering Woods and the jagged Peaks of Despair, once thrived under the benevolent rule of the Sunstone King. His power stemmed from the Orb of Aethelred, a mystical artifact radiating life and prosperity. But shadows stir. The Necromancer Malkor, banished centuries ago, has returned, corrupting the land with his undead legions. He seeks the Orb of Aethelred to plunge Eldoria into eternal darkness. A band of heroes must unite, brave treacherous landscapes, and confront Malkor before Eldoria is consumed by his malevolent reign."
}`

var roleResponse = `[
  {
    "name": "Lysandra",
    "description": "Lysandra, a wise Elven Mage, is Eldoria's last hope. Strength: 2, Dexterity: 3, Wisdom: 5. She draws upon ancient magic to protect her homeland from Gorgoth's encroaching darkness, harnessing elemental forces.",
    "effect": "Arcane Bolt: 4 MP, D8+3 damage; Shield: 3 MP, absorb 4 damage"
  }
]`

var eventResponse = `[
  {
    "name": "Whispering Woods Ambush",
    "description": "Combat event: Malkor's undead ambush the party within the Whispering Woods, seeking to halt their progress. Skeletal archers rain down poisoned arrows.",
    "effect": "Combat: Face 3 Skeletal Archers (HP 6, Attack D4 Poisoned) or lose 1d4 HP to poison."
  },
  {
    "name": "Ancient Elven Shrine",
    "description": "Plot event: An ancient Elven shrine, untouched by Malkor's corruption, offers guidance and forgotten lore to aid the heroes in their quest.",
    "effect": "Plot: Gain knowledge of Malkor's weakness. Advance one space on the Plot track."
  },
  {
    "name": "Dragon's Tooth Outpost",
    "description": "Combat event: Orcs loyal to Malkor control a strategic outpost in Dragon's Tooth. The heroes must reclaim it to secure a path.",
    "effect": "Combat: Face 5 Orc Warriors (HP 8, Attack D6) and an Orc Shaman (HP 12, Attack D4 + Magic)."
  },
  {
    "name": "Aethel's Echo",
    "description": "Plot event: The heroes find a fragment of the Orb of Aethel's power, resonating within an ancient ruin. It pulses with potent energy.",
    "effect": "Plot: Gain a temporary magical ability (+2 to any one roll) for the next three turns."
  },
  {
    "name": "Necromantic Ritual",
    "description": "Plot event: The heroes stumble upon a necromantic ritual site where Malkor is raising undead. They must disrupt the dark magic.",
    "effect": "Plot: Destroy the ritual. Reduce Malkor's army size (remove one minor enemy card from his forces)."
  },
  {
    "name": "Potion of Resistance",
    "description": "Item event: A hidden cache reveals a potent potion, offering temporary protection against Malkor's dark magic.",
    "effect": "Item: Potion of Resistance. Grants resistance to necrotic damage for 3 turns."
  }
]`

func main() {
	cards := []model.Card{}

	var story map[string]string
	if err := json.Unmarshal([]byte(storyResponse), &story); err != nil {
		fmt.Printf("Error unmarshalling AI response: %s\n", err)
		return
	}
	fmt.Println("Story generated successfully")
	fmt.Printf("Story: %s\n", story["story_background"])

	var role []cardResponse
	if err := json.Unmarshal([]byte(roleResponse), &role); err != nil {
		fmt.Printf("Error unmarshalling AI response: %s\n", err)
		return
	}

	fmt.Println("Role card generated successfully")

	fmt.Printf("Role: %+v\n", role)

	for _, r := range role {
		cards = append(cards, model.Card{
			GameID:      1, // Example GameID
			Type:        "role",
			Name:        r.Name,
			Description: r.Description,
			Effect:      r.Effect,
		})
	}
	fmt.Println("Card added to the game successfully")

	fmt.Printf("Generated Card: %+v\n", cards[0])

	var event []cardResponse
	if err := json.Unmarshal([]byte(eventResponse), &event); err != nil {
		fmt.Printf("Error unmarshalling AI response: %s\n", err)
		return
	}
	fmt.Println("Event card generated successfully")
	for _, e := range event {
		cards = append(cards, model.Card{
			GameID:      1, // Example GameID
			Type:        "event",
			Name:        e.Name,
			Description: e.Description,
			Effect:      e.Effect,
		})
	}
	fmt.Println("Event card added to the game successfully")
	for _, card := range cards {
		fmt.Printf("Card: %+v\n", card)
	}
}
