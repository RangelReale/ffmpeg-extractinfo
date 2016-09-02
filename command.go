package ffmpeg_info

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func (i *Info) execCmd(cmd *exec.Cmd) error {
	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
		//fmt.Printf("DONE CHANNEL EXIT\n")
	}()

	xlogout := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(cmdOutReader)
		scanner.Split(scanLines)
		for scanner.Scan() {
			//xlogout <- strings.TrimSpace(scanner.Text())
			xlogout <- strings.Trim(scanner.Text(), "\n\r")
		}
		//fmt.Printf("XLOGOUT CHANNEL EXIT\n")
	}()

	xlogerr := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(cmdErrReader)
		scanner.Split(scanLines)
		for scanner.Scan() {
			//xlogerr <- strings.TrimSpace(scanner.Text())
			xlogerr <- strings.Trim(scanner.Text(), "\n\r")
		}
		//fmt.Printf("XLOGERR CHANNEL EXIT\n")
	}()

	finished := false
	var execerr error

	t := time.NewTimer(i.FFMpegTimeout)
	for !finished {
		select {
		case <-t.C:
			{
				if err := cmd.Process.Kill(); err != nil {
					execerr = errors.New(fmt.Sprintf("timeout reached, but failed to kill: %s", err.Error()))
				}
				execerr = errors.New("process killed as timeout reached")
				finished = true
			}
		case l := <-xlogout:
			{
				if i.Event != nil {
					i.Event.OnStdout(l)
				}
			}
		case l := <-xlogerr:
			{
				if i.Event != nil {
					i.Event.OnStderr(l)
				}
			}
		case err := <-done:
			{
				if err != nil {
					execerr = errors.New(fmt.Sprintf("process done with error = %s", err.Error()))
				}
				finished = true
			}
		}
	}

	t.Stop()
	//close(done)

	return execerr
}

// Treats CRLF, CR not followed by LF and LF followed by anything as line endings
// https://groups.google.com/forum/#!topic/golang-nuts/cXX169-pNqw
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	var current, previous byte
	var offset int
loop:
	for advance, current = range data {
		switch {
		case current == '\n':
			break loop
		case previous == '\r':
			advance--
			break loop
		case current == '\r':
			previous = '\r'
		default:
			offset++
		}
	}
	token = data[:offset]
	advance++
	return
}
