package global

const (
	StoryPromptTemplate = "Generate a 100-word D&D-style fantasy story background for a board game. Include a setting, a central artifact, and a looming threat. Theme: %s, with json format: {\"story_background\": \"<story>\"}"

	RolePrompt = `Generate %s D&D-style characters for a board game based on story background: %s. Return a JSON object with:
		- "name": string (e.g., "Aragorn")
		- "description": string (50-word background, include profession like Warrior/Mage and attributes: Strength, Dexterity, Wisdom, range 1-5)
		- "effect": string (1-2 skills, e.g., "Fireball: 3 MP, D6+2 damage; Heal: 2 MP, restore 5 HP")
		Example:
		{
		"name": "Aragorn",
		"description": "Aragorn, a skilled Warrior. Strength: 4, Dexterity: 3, Wisdom: 2. A lone wanderer seeking an ancient artifact to defeat a dark lord.",
		"effect": "Sword Strike: D20+4 ≥ 15, D6+3 damage"
		}`

	EventPrompt = `Generate %s D&D-style board game event cards based on story background: %s. Return a JSON object with:
		- "name": string (e.g., "Dragon Attack")
		- "description": string (50-word description tied to the background)
		- "effect": string (e.g., "Combat: HP 10, Attack D6+1" or "Plot: Gain 1 Plot Point")
		Card type must be one of: combat, plot, item. Include type in description (e.g., "Combat event: ...").
		Example:
		{
		"name": "Dragon Attack",
		"description": "Combat event: A fire-breathing dragon assaults the village, demanding tribute. Heroes must fight to protect the innocent.",
		"effect": "Combat: HP 10, Attack D6+1"
		}`
)
