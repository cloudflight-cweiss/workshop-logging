package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"log/slog"
)

const (
	ApplicationName = "workshop-logapp"
)

func main() {
	lograteString := flag.String("lograte", "200", "Sets the lograte in milliseconds")
	lootSeedString := flag.String("seed", "1337", "Sets the initial loot seed")
	logType := flag.String("log", "text", "Sets the log output type. Allowed values are text or json")
	flag.Parse()

	switch *logType {
	case "text": // Allowed value but uses all default settings
		break
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	default:
		log.Fatalf("Unkown log output type '%s'", *logType)
	}

	slog.Info(fmt.Sprintf("Running %s ...", ApplicationName))

	lograteInt, err := strconv.Atoi(*lograteString)
	if err != nil {
		log.Fatalf("Could not parse --lograte argument: '%s'", *lograteString)
	}
	lograte := time.Duration(lograteInt) * time.Millisecond

	lootSeed, err := strconv.Atoi(*lootSeedString)
	if err != nil {
		log.Fatalf("Could not parse --lograte argument: '%s'", *lograteString)
	}
	var loot = NewLootTable(DefaultMessages, int64(lootSeed))

	slog.Info(fmt.Sprintf("Logging at %s rate with seed %d", lograte.String(), lootSeed))
	for {
		m := loot.Roll()
		slog.Log(context.Background(), m.Level, m.Message, "component", m.Component)
		time.Sleep(lograte)
	}
}
