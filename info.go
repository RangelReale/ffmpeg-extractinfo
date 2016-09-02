package ffmpeg_info

import (
	"time"
)

type Info struct {
	FFMpeg        *FFMpeg
	Event         Event
	ProbeTimeout  time.Duration
	FFMpegTimeout time.Duration
}

func NewInfo(FFMPeg *FFMpeg) *Info {
	return &Info{
		FFMpeg:        FFMPeg,
		ProbeTimeout:  2 * time.Second,
		FFMpegTimeout: 120 * time.Second,
	}
}
