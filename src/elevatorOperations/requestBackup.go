package elevatorOperations

import (
	"driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func TakeBackup() {
	requests := Fsm_requests()
	data, err := json.Marshal(requests)
	if err != nil {
		log.Println(Yellow, "Marshal conversion failed", White)
	}
	if err := ioutil.WriteFile("RequestBackup", data, 0644); err != nil {
		log.Println(Yellow, "ioutil.WriteFile() failed", White)
	}
}

func LoadBackup() {
	var requests [driver.N_FLOORS][driver.N_BUTTONS]bool
	if _, err := os.Stat("RequestBackup"); err == nil {
		log.Println("RequestBackup found, restoring ...")

		data, err := ioutil.ReadFile("RequestBackup")
		if err != nil {
			log.Println(Yellow, "Failed to read file from disk", White)
		}
		if err := json.Unmarshal(data, &requests); err != nil {
			log.Println(Yellow, "Marshal conversion failed", White)
			fmt.Println(err)
		}
		Fsm_printrequest(requests, "Restoring backup requests")
		for i := 0; i < 4; i++ {
			if requests[i][B_Inside] {
				Fsm_setRequest(i, B_Inside)
			}
		}
	}
}
