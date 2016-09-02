package ffmpeg_info

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type ProbeInfo_Format struct {
	Filename         string
	Nb_streams       int
	Nb_programs      int
	Format_name      string
	Format_long_name string
	Start_time       string
	Duration         string
	Size             string
	Bit_rate         string
	Probe_score      int
	Tags             map[string]string
}

type ProbeInfo_Stream struct {
	Index                int
	Codec_name           string
	Codec_long_name      string
	Profile              string
	Codec_type           string
	Codec_time_base      string
	Codec_tag_string     string
	Codec_tag            string
	Width                int
	Height               int
	Coded_width          int
	Coded_height         int
	Has_b_frames         int
	Sample_aspect_ratio  string
	Display_aspect_ratio string
	Pix_fmt              string
	Level                int
	Chroma_location      string
	Refs                 int
	Quarter_sample       string
	Divx_packed          string
	R_frame_rate         string
	Avg_frame_rate       string
	Time_base            string
	Start_pts            int
	Start_time           string
	Duration_ts          int
	Duration             string
	Bit_rate             string
	Nb_frames            string
	Tags                 map[string]string
}

type ProbeInfo struct {
	Format  *ProbeInfo_Format
	Streams []*ProbeInfo_Stream
}

func (i *Info) ProbeString(filename string) (string, error) {
	args := []string{
		"-i",
		filename,
		"-v",
		"quiet",
		"-print_format",
		"json",
		"-show_format",
		"-show_streams",
		"-show_error",
	}
	cmd := exec.Command(i.FFMpeg.exe_ffprobe, args...)

	var out bytes.Buffer
	var outerr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &outerr

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	var execerr error

	select {
	case <-time.After(i.ProbeTimeout):
		{
			if err := cmd.Process.Kill(); err != nil {
				execerr = errors.New(fmt.Sprintf("timeout reached, but failed to kill: %s", err.Error()))
			}
			execerr = errors.New("process killed as timeout reached")
		}
	case err := <-done:
		{
			if err != nil {
				execerr = errors.New(fmt.Sprintf("process done with error = %s", err.Error()))
			}
		}
	}

	if i.Event != nil {
		if out.Len() > 0 {
			//if len(out.String()) > 0 {
			i.Event.OnStdout(out.String())
		}
		if outerr.Len() > 0 {
			//if len(outerr.String()) > 0 {
			i.Event.OnStderr(outerr.String())
		}
	}

	if execerr != nil {
		return "", execerr
	}

	return out.String(), nil
}

func (i *Info) Probe(filename string) (*ProbeInfo, error) {
	str, err := i.ProbeString(filename)
	if err != nil {
		return nil, err
	}

	ret := &ProbeInfo{}
	err = json.Unmarshal([]byte(str), ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (i *Info) ProbeRaw(filename string) (map[string]interface{}, error) {
	str, err := i.ProbeString(filename)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal([]byte(str), &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
