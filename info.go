package ffmpeg_info

type Info struct {
	FFMpeg *FFMpeg
	Event  Event
}

func NewInfo(FFMPeg *FFMpeg) *Info {
	return &Info{
		FFMpeg: FFMPeg,
	}
}
