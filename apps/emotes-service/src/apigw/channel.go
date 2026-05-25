package apigw

type Channel struct {
	Id          string `json:"id"`
	TwitchId    string `json:"twitchId"`
	DisplayName string `json:"displayName"`
	ImageUrl    string `json:"imageUrl"`
	AddedAt     string `json:"addedAt"`
	UpdatedAt   string `json:"updatedAt"`
}
