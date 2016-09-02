package ffmpeg_info

type FFMpeg struct {
	exe_ffmpeg  string
	exe_ffprobe string
}

func NewFFMpeg() *FFMpeg {
	return &FFMpeg{
		exe_ffmpeg:  "ffmpeg",
		exe_ffprobe: "ffprobe",
	}
}

func NewFFMpegCustom(exe_ffmpeg string, exe_ffprobe string) *FFMpeg {
	return &FFMpeg{
		exe_ffmpeg:  exe_ffmpeg,
		exe_ffprobe: exe_ffprobe,
	}
}
