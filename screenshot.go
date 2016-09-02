package ffmpeg_info

import (
	"os/exec"
)

func (i *Info) Screenshot(filename string, outfilename string, time string) error {

	args := []string{
		"-y",
		"-ss",
		time,
		"-i",
		filename,
		"-vframes",
		"1",
		"-q:v",
		"2",
		outfilename,
	}
	cmd := exec.Command(i.FFMpeg.exe_ffmpeg, args...)

	err := i.execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}
