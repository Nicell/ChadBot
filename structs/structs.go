package structs

type Config struct {
	Discord string `json:"discord"`
	Youtube string `json:"youtube"`
}

type QueueItem struct {
	Title string
	URL   string
	User  string
}

type Guild struct {
	Queue []QueueItem
	Pause bool
	Skip  bool
}
