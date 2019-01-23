package music

import (
	"os"
	"os/exec"

	"github.com/rylio/ytdl"
)

func DlSng(vid *ytdl.VideoInfo, dcaPath string) (err error) {

	filePath := "downloads/" + vid.ID + ".webm"
	file, _ := os.Create(filePath)
	fileOut, _ := os.Create(dcaPath)
	defer file.Close()
	defer fileOut.Close()
	vid.Download(vid.Formats.Best("audenc")[0], file)

	cmd1 := exec.Command("ffmpeg", "-i", filePath, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	cmd2 := exec.Command("dca")

	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd2.Stdout = fileOut

	err = cmd2.Start()
	err = cmd1.Run()
	err = cmd2.Wait()
	if err != nil {
		return err
	}

	file.Close()
	os.Remove(filePath)

	return nil
}
