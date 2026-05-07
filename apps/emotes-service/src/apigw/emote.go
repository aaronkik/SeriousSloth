package apigw

import "emotes-service/src/adapters/secondary/event_store"

type EmoteImages struct {
	URL1X string `json:"url_1x"`
	URL2X string `json:"url_2x"`
	URL4X string `json:"url_4x"`
}

type Emote struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Format    []string    `json:"format"`
	Scale     []string    `json:"scale"`
	ThemeMode []string    `json:"theme_mode"`
	Images    EmoteImages `json:"images"`
}

func EmoteFromEvent(e event_store.EmoteServiceEventEmote) Emote {
	return Emote{
		ID:        e.ID,
		Name:      e.Name,
		Format:    e.Format,
		Scale:     e.Scale,
		ThemeMode: e.ThemeMode,
		Images: EmoteImages{
			URL1X: e.Images.URL1X,
			URL2X: e.Images.URL2X,
			URL4X: e.Images.URL4X,
		},
	}
}
