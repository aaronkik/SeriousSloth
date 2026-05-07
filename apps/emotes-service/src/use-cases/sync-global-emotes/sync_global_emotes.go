package syncglobalemotes

import (
	"context"
	"crypto/rand"
	"emotes-service/src/adapters/secondary/event_store"
	"emotes-service/src/adapters/secondary/twitch"
	"encoding/hex"
	"fmt"
	"log"
	"sort"
	"time"
)

type EmotesAggregate struct {
	State          map[string]struct{}
	LatestSequence int
}

func Aggregate(events []event_store.EmoteServiceEvent) *EmotesAggregate {
	agg := &EmotesAggregate{
		State: make(map[string]struct{}),
	}
	for _, event := range events {
		switch event.EventName {
		case "EmoteAdded":
			agg.State[event.EmoteId] = struct{}{}
		case "EmoteRemoved":
			delete(agg.State, event.EmoteId)
		}
		agg.LatestSequence = event.Sequence
	}
	return agg
}

func DecideSyncEvents(aggregateId string, agg *EmotesAggregate, globalEmotes []twitch.GlobalEmote) []event_store.EmoteServiceEvent {
	twitchState := make(map[string]twitch.GlobalEmote, len(globalEmotes))
	for _, emote := range globalEmotes {
		twitchState[emote.ID] = emote
	}

	var removedIds []string
	for id := range agg.State {
		if _, exists := twitchState[id]; !exists {
			removedIds = append(removedIds, id)
		}
	}
	sort.Strings(removedIds)

	var addedIds []string
	for id := range twitchState {
		if _, exists := agg.State[id]; !exists {
			addedIds = append(addedIds, id)
		}
	}
	sort.Strings(addedIds)

	currentSequence := agg.LatestSequence
	results := make([]event_store.EmoteServiceEvent, 0, len(removedIds)+len(addedIds))

	for _, id := range removedIds {
		currentSequence++
		results = append(results, createEmoteEvent(createEmoteEventInput{
			AggregateId: aggregateId,
			Sequence:    currentSequence,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
			EventName:   "EmoteRemoved",
			EmoteId:     id,
			Emote:       nil,
		}))
	}

	for _, id := range addedIds {
		currentSequence++
		src := twitchState[id]
		emote := &event_store.EmoteServiceEventEmote{
			Format:    src.Format,
			ID:        src.ID,
			Name:      src.Name,
			Scale:     src.Scale,
			ThemeMode: src.ThemeMode,
			Images: event_store.EmoteServiceEventEmoteImages{
				URL1X: src.Images.URL1X,
				URL2X: src.Images.URL2X,
				URL4X: src.Images.URL4X,
			},
		}

		results = append(results, createEmoteEvent(createEmoteEventInput{
			AggregateId: aggregateId,
			Sequence:    currentSequence,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
			EventName:   "EmoteAdded",
			EmoteId:     id,
			Emote:       emote,
		}))
	}

	return results
}

type createEmoteEventInput struct {
	AggregateId string
	Sequence    int
	CreatedAt   string
	EventName   string
	EmoteId     string
	Emote       *event_store.EmoteServiceEventEmote
}

func createEmoteEvent(in createEmoteEventInput) event_store.EmoteServiceEvent {
	return event_store.EmoteServiceEvent{
		PK:          in.AggregateId,
		SK:          fmt.Sprintf("SEQUENCE#%s", generateEventSequence(in.Sequence)),
		AggregateId: in.AggregateId,
		CreatedAt:   in.CreatedAt,
		Emote:       in.Emote,
		EmoteId:     in.EmoteId,
		EventName:   in.EventName,
		Id:          generateId(),
		Kind:        "EVENT",
		Sequence:    in.Sequence,
	}
}

func generateEventSequence(n int) string {
	return fmt.Sprintf("%07d", n)
}

func generateId() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("es_%s", hex.EncodeToString(bytes))
}

func Execute(ctx context.Context) error {
	token, err := twitch.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	emotes, err := twitch.GetGlobalEmotes(ctx, token)
	if err != nil {
		return err
	}

	events, err := event_store.LoadEvents(ctx, event_store.GlobalEmotesAggregateId)
	if err != nil {
		return err
	}

	aggregate := Aggregate(events)
	newEvents := DecideSyncEvents(event_store.GlobalEmotesAggregateId, aggregate, emotes)
	return event_store.AppendEvents(ctx, newEvents)
}
