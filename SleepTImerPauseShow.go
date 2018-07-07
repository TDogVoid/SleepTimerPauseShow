package main

import (
	"fmt"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/gen2brain/beeep"
	"github.com/micmonay/keybd_event"
)

var pauseAfter = 30
var startAlerts = 5

var lastInput = getLastInput()

var lastInputInfo struct {
	cbSize uint32
	dwTime uint32
}

func main() {
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
			fmt.Println("sleeping", i)
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
	if t == 1 {
		err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		if err != nil {
			panic(err)
		}
	}
}
