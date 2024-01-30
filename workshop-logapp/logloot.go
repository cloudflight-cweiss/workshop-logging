package main

import (
	"log"
	"log/slog"
	"math/rand"
	"slices"
)

type LogMessage struct {
	Level     slog.Level
	Message   string
	Component string
	Weight    int
}

var DefaultMessages = []LogMessage{
	{
		Level:     slog.LevelInfo,
		Message:   "This is an info log :D!",
		Component: "main",
		Weight:    10,
	},
	{
		Level:     slog.LevelInfo,
		Message:   "This is another info message, not really important",
		Component: "main",
		Weight:    10,
	},
	{
		Level:     slog.LevelInfo,
		Message:   "This is an info message from a subcomponent of the application",
		Component: "auth",
		Weight:    10,
	},
	{
		Level:     slog.LevelDebug,
		Message:   "This is a debug message, i should not be seen in production!",
		Component: "main",
		Weight:    10,
	},
	{
		Level:     slog.LevelWarn,
		Message:   "I am a warning, annoying but sometimes useful",
		Component: "auth",
		Weight:    10,
	},
	{
		Level:     slog.LevelError,
		Message:   "Regular error....somebody should really fix me...",
		Component: "calc",
		Weight:    10,
	},
	{
		Level:     slog.LevelError,
		Message:   "Error, error...this is a really rare error...fix me please!",
		Component: "main",
		Weight:    1,
	},
}

type LootTable struct {
	originalTable    []LogMessage
	replacementTable []LogMessage
	r                *rand.Rand
}

func NewLootTable(messageTable []LogMessage, seed int64) LootTable {
	return LootTable{
		originalTable:    slices.Clone(messageTable),
		replacementTable: slices.Clone(messageTable),
		r:                rand.New(rand.NewSource(seed)),
	}
}

func (t LootTable) resetLoot() {
	for i := range t.replacementTable {
		t.replacementTable[i].Weight = t.originalTable[i].Weight
	}
}

func (t LootTable) PrintLootProbability() {
	log.Fatalf("TBD: I should print loot probabilities!")
}

func (t LootTable) Roll() LogMessage {
	total := totalWeight(t.replacementTable)
	if total == 0 {
		t.resetLoot()
		total = totalWeight(t.replacementTable)
	}
	roll := t.r.Intn(total + 1)
	index := getWeightedItemIndex(t.replacementTable, roll)
	if index == -1 {
		panic("Selecting log message resulted in -1 index!")
	}
	t.replacementTable[index].Weight--
	return t.replacementTable[index]
}

func totalWeight(table []LogMessage) (total int) {
	for _, m := range table {
		total += m.Weight
	}
	return
}

func getWeightedItemIndex(table []LogMessage, roll int) int {
	rangeStart := 0
	rangeEnd := 0
	for i, m := range table {
		rangeEnd += m.Weight
		if roll >= rangeStart && roll <= rangeEnd {
			return i
		}
		rangeStart += m.Weight
	}
	return -1
}
