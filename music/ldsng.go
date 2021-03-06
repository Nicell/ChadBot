package music

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// LdSng loads the given video from a dca file and returns a buffer
func LdSng(videoID string) (buffer [][]byte, err error) {

	file, err := os.Open("library/" + videoID + ".dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return buffer, err
	}

	var opuslen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return buffer, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return buffer, err
		}

		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return buffer, err
		}

		buffer = append(buffer, InBuf)
	}
}
