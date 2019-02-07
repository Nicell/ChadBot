package music

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/rylio/ytdl"

	S "chadbot/structs"
)

// PlyNxt plays the next song in the given guild's queue
func PlyNxt(s *discordgo.Session, guildID, mChannelID string, vc *discordgo.VoiceConnection, voice *discordgo.Channel, guilds *map[string]*S.Guild) (err error) {
	vid, err := ytdl.GetVideoInfo((*guilds)[guildID].Queue[0].URL)
	if err != nil {
		s.ChannelMessageSend(mChannelID, "Failed to find the song!")
		return err
	}

	dcaPath := "library/" + vid.ID + ".dca"

	if _, err := os.Stat(dcaPath); os.IsNotExist(err) {

		err = DlSng(vid, dcaPath)
		if err != nil {
			return err
		}
	}

	buffer, err := LdSng(vid.ID)
	if err != nil {
		return err
	}

	s.ChannelMessageSend(mChannelID, "Now playing "+vid.Title+" in "+voice.Name)

	for _, buff := range buffer {
		for (*guilds)[guildID].Pause {
			if (*guilds)[guildID].Skip {
				break
			}
		}
		if (*guilds)[guildID].Skip {
			(*guilds)[guildID].Skip = false
			break
		}
		vc.OpusSend <- buff
	}

	(*guilds)[guildID].Queue = (*guilds)[guildID].Queue[1:]

	return nil
}
