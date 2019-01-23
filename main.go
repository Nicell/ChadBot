package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rylio/ytdl"

	H "chad/helpers"
	M "chad/music"
	S "chad/structs"
)

var guilds = make(map[string]*S.Guild, 0)

var cfg S.Config

func main() {

	err := H.LdCFG(&cfg)
	if err != nil {
		fmt.Println("Error parsing config: ", err)
		return
	}

	dg, err := discordgo.New("Bot " + cfg.Discord)

	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	os.Mkdir("downloads", 0777)
	os.Mkdir("library", 0777)

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(guildCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Println("ChadBot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	serverCount := strconv.Itoa(len(event.Guilds))
	s.UpdateListeningStatus(serverCount + " servers")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	tokens := strings.Split(m.Content, " ")

	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return
	}

	if tokens[0] == "!play" {

		if len(tokens) > 1 {
			for _, vs := range g.VoiceStates {
				if vs.UserID == m.Author.ID {
					err = handlePlay(s, g.ID, vs.ChannelID, m.ChannelID, tokens[1:], m.Author.ID)
					if err != nil {
						fmt.Println("Error playing sound:", err)
					}

					return
				}
			}

			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, I can't find you in any voice channel.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, you didn't enter any query `(Usage: !play {song})`.")
		}
	} else if tokens[0] == "!pause" {

		s.ChannelTyping(m.ChannelID)

		if _, in := guilds[g.ID]; in && len(guilds[g.ID].Queue) > 0 {

			guilds[g.ID].Pause = !guilds[g.ID].Pause

			if guilds[g.ID].Pause {
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, I just paused the music")
			} else {
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, I just unpaused the music")
			}
			return
		}

		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, there is no music playing")
	} else if tokens[0] == "!skip" {
		s.ChannelTyping(m.ChannelID)

		if _, in := guilds[g.ID]; in && len(guilds[g.ID].Queue) > 0 {

			guilds[g.ID].Skip = true

			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, I skipped "+guilds[g.ID].Queue[0].Title)
			return
		}

		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, there is no music playing")
	} else if tokens[0] == "!queue" {

		if _, in := guilds[g.ID]; in && len(guilds[g.ID].Queue) > 0 {

			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{},
				Color:  0x4caf50,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Now Playing: " + guilds[g.ID].Queue[0].Title,
						Value: "By <@" + guilds[g.ID].Queue[0].User + ">",
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
				Title:     "🎵 Music Queue 🎵",
			}

			for i := 1; i < len(guilds[g.ID].Queue); i++ {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:  strconv.Itoa(i) + ". " + guilds[g.ID].Queue[i].Title,
					Value: "By <@" + guilds[g.ID].Queue[i].User + ">",
				})
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}

		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, there is no music queued")
	} else if tokens[0] == "!remove" {

		if len(tokens) > 1 {

			qLen := len(guilds[g.ID].Queue)

			if _, in := guilds[g.ID]; in && qLen > 1 {

				pos, err := strconv.Atoi(tokens[1])
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, that isn't a valid position `(Usage: !remove {song position})`")
					return
				}

				if pos >= qLen || pos <= 0 {
					s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, position "+strconv.Itoa(pos)+" is out of bounds")
					return
				}

				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, I removed "+guilds[g.ID].Queue[pos].Title+" from the queue")
				guilds[g.ID].Queue = append(guilds[g.ID].Queue[:pos], guilds[g.ID].Queue[pos+1:]...)
				return
			}

			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, there is no music queued")
			return
		}

		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">, missing position to remove `(Usage: !remove {song position})`")
	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	serverCount := strconv.Itoa(len(s.State.Guilds))
	s.UpdateListeningStatus(serverCount + " servers")
}

func handlePlay(s *discordgo.Session, guildID, vcChannelID string, mChannelID string, arguments []string, authorID string) (err error) {

	s.ChannelTyping(mChannelID)

	link := arguments[0]

	if !H.ValidURL(link) {
		link, err = H.YTsearch(strings.Join(arguments, " "), cfg.Youtube)
		if err != nil {
			s.ChannelMessageSend(mChannelID, "<@"+authorID+">, failed to find the song!")
			return err
		}
	}

	vid, err := ytdl.GetVideoInfo(link)
	if err != nil {
		s.ChannelMessageSend(mChannelID, "<@"+authorID+">, failed to find the song!")
		return err
	}

	if _, in := guilds[guildID]; !in {
		guilds[guildID] = &S.Guild{
			Queue: make([]S.QueueItem, 0),
			Pause: false,
			Skip:  false,
		}
	}

	guilds[guildID].Queue = append(guilds[guildID].Queue, S.QueueItem{Title: vid.Title, URL: link, User: authorID})

	s.ChannelMessageSend(mChannelID, "<@"+authorID+">, I added "+vid.Title+" to the queue")

	if len(guilds[guildID].Queue) == 1 {
		M.PlyQ(s, guildID, vcChannelID, mChannelID, &guilds)
	}

	return nil
}
