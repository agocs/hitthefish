package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type item struct {
	name       string
	desc       string
	isHittable bool
	isHit      bool
}

type environment struct {
	name  string
	exits map[string]environment
	inv   inventory
}

func (e environment) String() string {
	outs := e.name
	outs += "\n"
	outs += "Exits are: "
	for k, _ := range e.exits {
		outs += k + " "
	}
	outs += "\n"
	if len(e.inv) > 0 {
		outs += "Room contains: \n"
		outs += e.inv.String()
	}
	return outs
}
func (e *environment) setExit(dir string, x environment) {
	e.exits[dir] = x
}

type inventory map[string]item

func (i inventory) String() string {
	outString := ""
	for k, _ := range i {
		outString += fmt.Sprint(k)
	}
	return outString
}

var personalInventory inventory = make(map[string]item)

var currentEnv = firstRoom

func main() {
	for {
		doEnvironment(currentEnv)
		if fishRoom.inv[theFish.name].isHit {
			fmt.Println("You win!")
			os.Exit(0)
		}
		if _, ok := fishRoom.inv[theFish.name]; !ok {
			fmt.Println("You lose.")
			os.Exit(1)
		}
		fmt.Println()
	}
}

var firstRoom = environment{
	name:  "The first room",
	exits: make(map[string]environment),
	inv: inventory{"fish whacker": item{
		name: "fish whacker",
		desc: "a fine stick with which to hit the fish",
	}},
}

var theFish = item{
	name:       "fish",
	desc:       "a fine hittable fish",
	isHittable: true,
}

var fishRoom = environment{
	name:  "The fish room",
	exits: make(map[string]environment),
	inv:   inventory{"fish": theFish},
}

func init() {
	firstRoom.setExit("north", fishRoom)
	fishRoom.setExit("south", firstRoom)
}

func doEnvironment(env environment) {
	fmt.Println(env)
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	sentence, _ := buf.ReadBytes('\n')
	doTask(string(sentence[:len(sentence)-1]))

}

var verbs = map[string]string{
	"north":     "north",
	"south":     "south",
	"east":      "east",
	"west":      "west",
	"inventory": "inventory",
	"look":      "look",
	"take":      "take",
	"hit":       "hit",
	"quit":      "quit",
}

func doTask(sentence string) {
	words := strings.Split(sentence, " ")
	fmt.Println(words)
	switch words[0] {
	case verbs["north"], verbs["south"], verbs["east"], verbs["west"]:
		makeMove(words[0])
	case verbs["take"]:
		takeItem(strings.Join(words[1:], " "))
	case verbs["inventory"]:
		fmt.Println(personalInventory)
	case verbs["quit"]:
		os.Exit(0)
	case verbs["look"]:
		describeItem(strings.Join(words[1:], " "))
	case verbs["hit"]:
		hitItem(strings.Join(words[1:], ""))
	default:
		fmt.Println("I don't know what you mean.")
		fmt.Println("Valid verbs are: ")
		for k, _ := range verbs {
			fmt.Println(k)
		}
	}
}

func makeMove(direction string) {
	newEnv, ok := currentEnv.exits[direction]
	if !ok {
		fmt.Printf("There's no exit to the %s\n", direction)
		return
	}
	currentEnv = newEnv
}

func takeItem(name string) {
	item, ok := currentEnv.inv[name]
	if !ok {
		fmt.Printf("I don't see a %s here\n", name)
		return
	}
	personalInventory[name] = item
	delete(currentEnv.inv, name)
}

func describeItem(name string) {
	if i, ok := personalInventory[name]; ok {
		fmt.Println(i.desc)
		fmt.Printf("You are holding the %s\n", i.name)
		return
	}
	if i, ok := currentEnv.inv[name]; ok {
		fmt.Println(i.desc)
		return
	}
	fmt.Printf("I don't see the %s here.\n", name)

}

func hitItem(name string) {
	if _, ok := personalInventory["fish whacker"]; !ok {
		fmt.Println("With what?")
		return
	}

	if i, ok := currentEnv.inv[name]; ok {
		if !i.isHittable {
			fmt.Printf("You can't hit the %s\n", name)
			return
		}
		i.isHit = true
		currentEnv.inv[name] = i
		fmt.Println("good hit")
		return
	}
	fmt.Println("I don't see that here!")
}
