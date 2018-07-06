package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/micmonay/keybd_event"
)

var lastInput = getLastInput()

var lastInputInfo struct {
	cbSize uint32
	dwTime uint32
}

func main() {
	timer(30) // in minutes
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
