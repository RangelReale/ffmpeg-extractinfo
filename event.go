package ffmpeg_info

type Event interface {
	OnStdout(text string)
	OnStderr(text string)
}
