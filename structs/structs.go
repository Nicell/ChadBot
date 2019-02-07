package structs

// Config holds a config's settings
type Config struct {
	Discord string
	Youtube string
}

// QueueItem holds a song's information
type QueueItem struct {
	Title string
	URL   string
	User  string
}

// Guild holds a guild's state
type Guild struct {
	Queue []QueueItem
	Pause bool
	Skip  bool
}
