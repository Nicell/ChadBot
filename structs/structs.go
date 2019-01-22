package structs

type Config struct {
	Discord string `json:"discord"`
	Youtube string `json:"youtube"`
}

type QueueItem struct {
	Title string
	URL   string
}

type Guild struct {
	Queue []QueueItem
	Pause bool
	Skip  bool
}
