package player

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

var (
	cmd     *exec.Cmd
	inPipe  io.WriteCloser
	outPipe *bufio.Reader
	ch      = make(chan int)
)

func PlayAndWait(audio string) error {
	if err := Stop(); err != nil {
		return err
	}
	if audio == "" {
		return nil
	}
	loadCmd := "l " + audio + "\n"
	if _, err := inPipe.Write([]byte(loadCmd)); err != nil {
		return err
	}

	result := 1
	for result == 1 { // wait for stop
		result = <-ch
		fmt.Println("play: ", result)
	}
	return nil
}

func PauseOrResume() (err error) {
	loadCmd := "p\n"
	_, err = inPipe.Write([]byte(loadCmd))
	return
}

func Next() (err error) {
	loadCmd := "s\n"
	_, err = inPipe.Write([]byte(loadCmd))
	return
}

func Stop() (err error) {
	loadCmd := "s\n"
	if _, err = inPipe.Write([]byte(loadCmd)); err == nil {
		<-ch
	}
	return
}

func process() {
	for {
		_line, _, err := outPipe.ReadLine()
		line := string(_line)
		if err != nil {
			if err == io.EOF {
				ch <- 0
				fmt.Println("process EOF: ", 0)
			} else {
				ch <- -1
				fmt.Println("process else: ", 1)
			}
			continue
		}
		if strings.HasPrefix(line, "@P 0") {
			ch <- 0
			fmt.Println("process @P 0: ", 0)
		} else if strings.HasPrefix(line, "@P 1") {
			ch <- 1
			fmt.Println("process @P 1: ", 1)
		}
	}
}

func StartAndWait() {
	cmd = exec.Command("sudo", "mpg123", "-R")
	var err error
	inPipe, err = cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	_outPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	outPipe = bufio.NewReader(_outPipe)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("Waiting for command to finish...")
	go process()
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	fmt.Printf("mpg123 stoped.\n")
}
