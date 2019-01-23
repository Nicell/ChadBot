package music

import (
	"github.com/bwmarrin/discordgo"

	S "chad/structs"
)

// PlyQ plays the queue of the given guild
func PlyQ(s *discordgo.Session, guildID, vcChannelID string, mChannelID string, guilds *map[string]*S.Guild) (err error) {

	vc, err := s.ChannelVoiceJoin(guildID, vcChannelID, false, true)
	if err != nil {
		return err
	}

	voice, err := s.Channel(vcChannelID)
	if err != nil {
		return err
	}

	vc.Speaking(true)

	for len((*guilds)[guildID].Queue) > 0 {
		err = PlyNxt(s, guildID, mChannelID, vc, voice, &*guilds)
		if err != nil {
			return err
		}
	}

	vc.Speaking(false)

	vc.Disconnect()

	return nil
}
