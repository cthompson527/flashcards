package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/term"
)

// Returns either an ascii code, or (if input is an arrow) a Javascript key code.
// code taken from: https://github.com/paulrademacher/climenu/blob/a1afbb4e378bf580e7d6bddd826e44e8f64347a1/getchar.go#L6
func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}

type definition struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

func getDefinitions() []definition {
	raw, err := ioutil.ReadFile("./terms.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var d []definition
	json.Unmarshal(raw, &d)
	return d
}

func getUserInput(valid []int) int {
	for {
		ascii, _, err := getChar()
		checkForError(err)

		for _, v := range valid {
			if v == ascii {
				return ascii
			}
		}

		if ascii == 'q' {
			os.Exit(0)
		}
	}
}

func checkForError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func clearScreen() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	clearScreen()
	def := getDefinitions()
	length := len(def)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("*** Welcome to Flashcards ***")
	fmt.Println("Press ? for instructions, or b to begin...")

	var ascii int

	ascii = getUserInput([]int{'?', 'b'})

	if ascii == '?' {
		fmt.Println("1. When only the term is displayed, press spacebar to view definition")
		fmt.Println("2. After definition is displayed, press spacebar to go again")
		fmt.Println("3. At any time, press q to quit")
		fmt.Println("\nPress b to begin...")

		ascii = getUserInput([]int{'b'})
	}

	for {
		elem := def[r.Intn(length)]
		fmt.Printf("%s\t-\t", elem.Term)

		ascii = getUserInput([]int{32})
		fmt.Printf("%s\n", elem.Definition)

		ascii = getUserInput([]int{32})
	}
}
