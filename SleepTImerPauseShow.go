package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/gen2brain/beeep"
	"github.com/micmonay/keybd_event"
)

var pauseAfter = 30 //Default Pause After Time in minutes
var startAlerts = 5 //Default Start Alert Time in minutes
var beepAlert = true

var lastInput = getLastInput()

var lastInputInfo struct {
	cbSize uint32
	dwTime uint32
}

func main() {
	getTimeForTimer()
	getAlertTime()
	getBeepAlert()
	timer(pauseAfter) // in minutes

}

func pauseShow() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(keybd_event.VK_SPACE)
	err = kb.Launching()
	if err != nil {
		panic(err)
	}
}

func wasThereInput() bool {
	if getLastInput() != lastInput {
		lastInput = getLastInput()
		return true
	}
	return false
}

//timer waits t number of minutes
func timer(t int) {
	i := 0
	for {
		isInput := wasThereInput()
		if !isInput {
			fmt.Printf("Been inactive for %v mins\n", i)
			alert(i)
			time.Sleep(1 * time.Minute)
			i++
		} else if isInput {
			i = 0 //restart timer

		}

		if i >= t {
			//pause show
			fmt.Println("pause")
			pauseShow()
			return
		}

	}

}

func alert(idleSince int) {
	if pauseAfter-idleSince <= startAlerts {
		notifyTimeBefore(pauseAfter - idleSince)
	}
}

// gets the last time user input in ms since system started
func getLastInput() uint32 {
	//From https://stackoverflow.com/questions/22949444/using-golang-to-get-windows-idle-time-getlastinputinfo-or-similar
	lastInputInfo.cbSize = uint32(unsafe.Sizeof(lastInputInfo))

	user32 := syscall.MustLoadDLL("user32.dll")                 // or NewLazyDLL() to defer loading
	getLastInputInfo := user32.MustFindProc("GetLastInputInfo") // or NewProc() if you used NewLazyDLL()
	// or you can handle the errors in the above if you want to provide some alternative
	r1, _, err := getLastInputInfo.Call(
		uintptr(unsafe.Pointer(&lastInputInfo)))
	// err will always be non-nil; you need to check r1 (the return value)
	if r1 == 0 { // in this case
		panic("error getting last input info: " + err.Error())
	}
	return lastInputInfo.dwTime
}

func notifyTimeBefore(t int) {
	fmt.Println("Notify")
	err := beeep.Notify("Sleep Timer", "Pausing your show in "+strconv.Itoa(t)+" mins\nMove your mouse or press a button to reset timer", "assets/information.png")
	if err != nil {
		panic(err)
	}
	if t == 1 && beepAlert {
		err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		if err != nil {
			panic(err)
		}
	}
}

func getTimeForTimer() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("How long do you want to set the timer in minutes? (default: %v) ", pauseAfter)
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\r\n")
	if text != "" {
		i, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Invalid Number: ", err)
			getTimeForTimer()
			return
		}
		pauseAfter = i
	}
	fmt.Printf("Will pause after %v mins of inactivity\n", pauseAfter)
}

func getAlertTime() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("At what time do you want alerts to start in mins? (default: %v) ", startAlerts)
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\r\n")
	if text != "" {
		i, err := strconv.Atoi(text)
		if err != nil {
			fmt.Println("Invalid Number: ", err)
			getTimeForTimer()
			return
		}
		startAlerts = i
	}

	fmt.Printf("Will start alerts %v mins before pause\n", startAlerts)
}

func getBeepAlert() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Do you wish to have a beep alert on last notification? y/n (default: y) ")
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\r\n")
	text = strings.ToLower(text)
	switch text {
	case "y":
		beepAlert = true
	case "n":
		beepAlert = false
	case "":
		beepAlert = true
	default:
		fmt.Println("Invalid y/n")
		getBeepAlert()
		return
	}
	if beepAlert {
		fmt.Println("Will beep on last notification")
	} else {
		fmt.Println("Won't beep")
	}
}
