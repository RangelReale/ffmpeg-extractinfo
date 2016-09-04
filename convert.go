package ffmpeg_info

import (
	"os/exec"
	"strconv"
	"strings"
)

type Convert_ScaleMode int

const (
	CScaleMode_LetterBox Convert_ScaleMode = 1
	CScaleMode_Fit                         = 2
	CScaleMode_Exact                       = 3
)

func (i *Info) Convert(filename string, outfilename string, width int, height int, scalemode Convert_ScaleMode,
	moreargs []string) error {
	sscale := ""
	if scalemode == CScaleMode_LetterBox {
		// letterbox
		// http://superuser.com/questions/547296/resizing-videos-with-ffmpeg-avconv-to-fit-into-static-sized-player
		sscale = `scale=iw*min($width/iw\,$height/ih):ih*min($width/iw\,$height/ih), pad=$width:$height:($width-iw*min($width/iw\,$height/ih))/2:($height-ih*min($width/iw\,$height/ih))/2`
	} else if scalemode == CScaleMode_Fit {
		// https://trac.ffmpeg.org/wiki/Scaling%20(resizing)%20with%20ffmpeg
		sscale = `scale=w=$width:h=$height:force_original_aspect_ratio=decrease`
	} else {
		sscale = `scale=w=$width:h=$height,setsar=1:1`
	}
	sscale = strings.Replace(sscale, "$width", strconv.FormatInt(int64(width), 10), -1)
	sscale = strings.Replace(sscale, "$height", strconv.FormatInt(int64(height), 10), -1)

	args := []string{
		"-y",
		//"-hide_banner",
		"-i",
		filename,
		"-filter:v",
		sscale,
	}
	if moreargs != nil {
		args = append(args, moreargs...)
	}
	args = append(args, outfilename)
	cmd := exec.Command(i.FFMpeg.exe_ffmpeg, args...)

	err := i.execCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}
