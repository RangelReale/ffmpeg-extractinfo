package ffmpeg_info

type ProbeInfo_Format struct {
}

type ProbeInfo_Stream struct {
}

type ProbeInfo struct {
	Format  *ProbeInfo_Format
	Streams []*ProbeInfo_Stream
}

func (i *Info) Probe(filename string) (*ProbeInfo, error) {
	return nil, nil
}
