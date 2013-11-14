package player

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

var (
	inPipe   io.WriteCloser
	outPipe  *bufio.Reader
	statusCh = make(chan int, 1) //0: stop, 1: pausing
)

/* would be block if can not play a new song */
func ReadyToPlay() {
	for {
		if status := <-statusCh; status == 0 {
			break
		}
		fmt.Println("wait for status 0")
	}
}

/** make player can be ready to play again */
func ResetToPlay() {
	statusCh <- 0
}

/* not thread safe */
func IsStop() bool {
	select {
	case status := <-statusCh:
		stoped := status == 0
		statusCh <- status
		return stoped
	default:
		return false
	}
	return false
}

func Play(audio string) error {
	return sendCmd("l " + audio + "\n")
}

func PauseOrResume() (err error) {
	return sendCmd("p\n")
}

func Next() (err error) {
	return sendCmd("s\n")
}

func Stop() (err error) {
	return sendCmd("s\n")
}

func sendCmd(cmd string) (err error) {
	_, err = inPipe.Write([]byte(cmd))
	return
}

func process() {
	for {
		_line, _, err := outPipe.ReadLine()
		line := string(_line)
		if err != nil {
			panic(err)
		}
		if strings.HasPrefix(line, "@P 0") { //music is stoped
			statusCh <- 0
			fmt.Println("process @P 0: ", 0)
		} else if strings.HasPrefix(line, "@P 1") { //music is pausing
			statusCh <- 1
			fmt.Println("process @P 1: ", 1)
		}
	}
}

func StartAndWait() {
	cmd := exec.Command("sudo", "mpg123", "-R")
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
