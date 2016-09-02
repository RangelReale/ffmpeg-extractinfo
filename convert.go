package ffmpeg_info

import (
	"os/exec"
	"strconv"
	"strings"
)

func (i *Info) Convert(filename string, outfilename string, width int, height int) error {
	// letterbox
	// http://superuser.com/questions/547296/resizing-videos-with-ffmpeg-avconv-to-fit-into-static-sized-player
	sscale := `scale=iw*min($width/iw\,$height/ih):ih*min($width/iw\,$height/ih), pad=$width:$height:($width-iw*min($width/iw\,$height/ih))/2:($height-ih*min($width/iw\,$height/ih))/2`
	sscale = strings.Replace(sscale, "$width", strconv.FormatInt(int64(width), 10), -1)
	sscale = strings.Replace(sscale, "$height", strconv.FormatInt(int64(height), 10), -1)

	args := []string{
		"-y",
		//"-hide_banner",
		"-i",
		filename,
		"-filter:v",
		sscale,
		outfilename,
	}
	cmd := exec.Command(i.FFMpeg.exe_ffmpeg, args...)

	err := i.execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}
