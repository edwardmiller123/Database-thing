package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
)

type person struct {
	Name      string `json:"Name"`
	Age       int    `json:"Age"`
	FavAnimal string `'json:"FavAnimal"`
	IdNumber  int    `json:"IdNumber"`
}

func generateId(arr []person) int {
	Id := rand.Intn(100000)
	for m := 0; m < len(arr); m++ {

		if arr[m].IdNumber == Id {
			Id += 5
			generateId(arr)
		} else {

		}

	}
	return Id
}

func main() {

	fmt.Println("Database starting...")
	// connect to couchbase
	bucketName := "People"
	username := "Administrator"
	password := "password"

	cluster, err := gocb.Connect("couchbase://localhost", gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket(bucketName)
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	var selector, password2 string
	dataBase := []person{}

	time.Sleep(time.Second)

	defer fmt.Println("Data base shutting down.")

	for selector != "end" {

		for selector != "1" && selector != "2" && selector != "3" && selector != "end" {
			fmt.Printf("Press 1 to register new person.\nPress 2 to look up existing Person\nPress 3 for admin options.\n")
			fmt.Scan(&selector)
		}
		if selector == "1" { // register new person
			var newPerson person
			fmt.Println("Enter name, age and favourite animal")

			fmt.Scan(&newPerson.Name)
			fmt.Scan(&newPerson.Age)
			fmt.Scan(&newPerson.FavAnimal)

			newPerson.IdNumber = generateId(dataBase)

			dataBase = append(dataBase, newPerson)

			fmt.Println("Person added.")
			selector = "0"

		} else if selector == "2" { //search for existing
			var selector2 string

			for selector2 != "1" && selector2 != "2" && selector2 != "end" {
				fmt.Printf("To search by name or favourite animal press 1.\nTo search by age press 2.\n")
				fmt.Scan(&selector2)
				if selector2 == "end" {
					selector = "end"
				}
			}

			if selector2 == "1" {
				var search string
				fmt.Println("Enter name or animal:")
				fmt.Scan(&search)

				for i := 0; i < len(dataBase); i++ {

					if dataBase[i].Name == search || dataBase[i].FavAnimal == search {
						fmt.Println(dataBase[i])
					}
				}
				selector2 = "0"
				selector = "0"

			} else if selector2 == "2" {
				var searchAge int
				fmt.Println("Enter age:")
				fmt.Scan(&searchAge)

				for j := 0; j < len(dataBase); j++ {

					if dataBase[j].Age == searchAge {
						fmt.Println(dataBase[j])
					}
				}
			}
		} else if selector == "3" { //login

			tries := 0

			for password2 != "badpassword123" {
				fmt.Printf("Enter password: ")
				fmt.Scan(&password2)

				if tries == 2 {
					fmt.Println("One attempt remaining!")
				} else if tries == 3 {
					break
				}

				tries++
			}
			var selector3 string

			if password2 != "badpassword123" {
				selector = "0"
				selector3 = "return"
			}

			for selector3 != "1" && selector3 != "2" && selector3 != "return" && selector3 != "end" {
				fmt.Println(dataBase)
				fmt.Printf("To edit existing entry, press 1.\nTo return, type \"return\".\n")
				fmt.Scan(&selector3)

			}
			if selector3 == "end" {
				selector = "end"

			} else if selector3 == "return" {
				selector = "0"

			} else if selector3 == "1" { // edit entry
				var idToChange int
				var propertyToChange, newValue string

				fmt.Println("Enter Id number of entry to be modified:")
				fmt.Scan(&idToChange)

				for n := 0; n < len(dataBase); n++ {

					if dataBase[n].IdNumber == idToChange {

						fmt.Println("Enter property to be modified:")
						//This doesnt work for multiple words e.g "favourite animal"
						//Fix this!
						fmt.Scan(&propertyToChange)

						fmt.Println("Enter new value:")
						fmt.Scan(&newValue)

						if propertyToChange == "name" {
							dataBase[n].Name = newValue

						} else if propertyToChange == "age" {
							newValueInt, err := strconv.Atoi(newValue)

							if err != nil {
								panic(err)
							}

							dataBase[n].Age = newValueInt

						} else if propertyToChange == "favourite animal" {
							dataBase[n].FavAnimal = newValue
						}
					}
				}

				selector = "3"
				selector3 = "0"

			}

		}

	}
}
