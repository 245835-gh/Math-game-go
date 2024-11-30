package main

import (
	"Game1/domain"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	totalPoints      int = 100
	pointPerQuestion int = 100
)

var id uint64 = 1

func main() {
	users := getUsers()
	for _, user := range users {
		if user.Id >= id {
			id = user.Id + 1
		}
	}
	fmt.Println("Welcome to game")

	for {
		menu()
		point := ""
		fmt.Scan(&point)
		switch point {
		case "1":
			user := play()
			users = getUsers()
			users = append(users, user)
			sortAndSave(users)
		case "2":
			users = getUsers()
			for _, user := range users {
				fmt.Printf("Id: %v Name: %s Time: %v \n", user.Id, user.Name, user.TimeSpent)
			}
		case "3":
			return
		default:
			fmt.Println("change number 1,2 or 3")
		}
	}

}

func menu() {
	println("1. Startgame")
	println("2. Result")
	println("3. Exit")
}

func play() domain.User {
	for i := 5; i >= 1; i-- {
		fmt.Printf("Game start in: %v\n", i)
		time.Sleep(1 * time.Second)
	}

	startTime := time.Now()
	myPoint := 0
	for myPoint < totalPoints {
		x, y := rand.Intn(100), rand.Intn(100)
		fmt.Printf("%v + %v =", x, y)
		ans := ""
		fmt.Scan(&ans)
		ansInt, err := strconv.Atoi(ans)
		if err != nil {
			fmt.Println("Please write integer number")
		} else {
			if ansInt == x+y {
				myPoint += pointPerQuestion
				fmt.Println("My points:", myPoint)
				fmt.Printf("point end: %v\n", totalPoints-myPoint)
			} else {
				fmt.Println("Try again")
			}
		}
	}
	endTime := time.Now()
	timeSpent := endTime.Sub(startTime)
	fmt.Println("Game over you win (time):", timeSpent)
	fmt.Println("Write your name: ")

	name := ""
	fmt.Scan(&name)

	user := domain.User{
		Id:        id,
		Name:      name,
		TimeSpent: timeSpent,
	}
	id++
	return user
}

func sortAndSave(users []domain.User) {
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].TimeSpent < users[j].TimeSpent
	})

	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("sortAndSave -> os.OpenFile: %s\n", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(users)
	if err != nil {
		fmt.Printf("sortAndSave -> os.OpenFile: %s\n", err)
		return
	}
}

func getUsers() []domain.User {
	file, err := os.Open("users.json")
	if err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create("users.json")
			if err != nil {
				fmt.Printf("getUsers -> os.Create: %s\n", err)
				return nil
			}
			return nil
		}
		fmt.Printf("getUsers -> os.Open: %s\n", err)
		return nil
	}
	var users []domain.User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		fmt.Printf("getUsers -> decoder.Decode: %s\n", err)
		return nil
	}

	return users
}
