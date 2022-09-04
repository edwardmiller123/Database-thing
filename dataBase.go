package main

import (
	"fmt"
	"math/rand"
	"time"
)

type person struct {
	name      string
	age       int
	favAnimal string
	idNumber  int
}

func generateId(arr []person) int {
	Id := rand.Intn(5)
	for m := 0; m < len(arr); m++ {

		if arr[m].idNumber == Id {
			Id += 1
			generateId(arr)
		} else {

		}

	}
	return Id
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var selector string
	dataBase := []person{}
	for selector != "end" {
		for selector != "1" && selector != "2" && selector != "3" && selector != "end" {
			fmt.Printf("Press 1 to register new person.\nPress 2 to look up existing Person\nPress 3 to edit existing entry.\n")
			fmt.Scan(&selector)
		}
		if selector == "1" {
			var newPerson person
			fmt.Println("Enter name, age and favourite animal")

			fmt.Scan(&newPerson.name)
			fmt.Scan(&newPerson.age)
			fmt.Scan(&newPerson.favAnimal)

			newPerson.idNumber = generateId(dataBase)

			dataBase = append(dataBase, newPerson)

			fmt.Println("Person added.")
			selector = "0"

		} else if selector == "2" {
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
					if dataBase[i].name == search || dataBase[i].favAnimal == search {
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
					if dataBase[j].age == searchAge {
						fmt.Println(dataBase[j])
					}
				}
			}
		} else if selector == "3" {
			fmt.Println(dataBase)
			selector = "0"
		}

	}
}
