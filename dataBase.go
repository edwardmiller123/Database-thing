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
	FavAnimal string `'json:"FavouriteAnimal"`
	IdNumber  int    `json:"IdNumber"`
}

type loginCredentials struct {
	UserName     string `json:"UserName"`
	UserPassword string `json:"UserPassword"`
}

func generateId(arr []person) int {
	Id := rand.Intn(100000)
	for m := 0; m < len(arr); m++ {

		if arr[m].IdNumber == Id {
			Id += 5
			generateId(arr)
		}

	}
	return Id
}

func main() {

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

	col := bucket.Scope("_default").Collection("_default")
	col2 := bucket.Scope("_default").Collection("loginCredentials")

	rand.Seed(time.Now().UnixNano())

	var selector, password2 string
	var dataBase []person
	var admins []loginCredentials
	var skipLogin bool

	queryResult1, err := cluster.Query(fmt.Sprintf("select Name, Age, FavAnimal, IdNumber from `%s`._default._default", bucketName), &gocb.QueryOptions{})

	if err != nil {
		log.Fatal(err)
	}

	queryResult2, err := cluster.Query(fmt.Sprintf("select UserName, UserPassword from `%s`._default.loginCredentials", bucketName), &gocb.QueryOptions{})

	if err != nil {
		log.Fatal(err)
	}

	//Appends people to local variable
	for queryResult1.Next() {
		var result person
		err := queryResult1.Row(&result)
		if err != nil {
			log.Fatal(err)
		}
		dataBase = append(dataBase, result)
	}

	//Appends admins to local variable
	for queryResult2.Next() {
		var result2 loginCredentials
		err := queryResult2.Row(&result2)
		if err != nil {
			log.Fatal(err)
		}
		admins = append(admins, result2)
	}

	fmt.Printf("Database Active\n \n")

	defer fmt.Println("Data base shutting down.")

	fmt.Printf("Currently Registered: %d\n--------------------------\nNumber of Admins: %d\n--------------------------\n", len(dataBase), len(admins))

	for selector != "end" {

		for selector != "1" && selector != "2" && selector != "3" && selector != "end" {
			fmt.Printf("\nPress 1 to register new person.\n//\nPress 2 to look up existing Person\n//\nPress 3 for admin options.\n--------------------------\n")
			fmt.Scan(&selector)
		}
		if selector == "1" { // register new person
			var newPerson person
			fmt.Println("Enter name, age and favourite animal")

			fmt.Scan(&newPerson.Name)
			fmt.Scan(&newPerson.Age)
			//doesnt work for multiple words
			//fix this
			fmt.Scan(&newPerson.FavAnimal)

			newPerson.IdNumber = generateId(dataBase)

			dataBase = append(dataBase, newPerson)

			//Update couchbase
			_, err = col.Upsert(fmt.Sprintf("%d", newPerson.IdNumber),
				person{
					Name:      newPerson.Name,
					Age:       newPerson.Age,
					FavAnimal: newPerson.FavAnimal,
					IdNumber:  newPerson.IdNumber,
				}, nil)
			if err != nil {
				log.Fatal(err)
			}

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
				selector2 = "0"
			}
		} else if selector == "3" { //login
			var userCheck, selector3 string
			tries := 0

			if !skipLogin {
				fmt.Println("Enter user Name:")
				fmt.Scan(&userCheck)
			}

			for p := 0; p < len(admins); p++ {

				if userCheck == admins[p].UserName && !skipLogin {

					for password2 != admins[p].UserPassword {
						fmt.Printf("Enter password: ")
						fmt.Scan(&password2)

						if tries == 2 {
							fmt.Println("One attempt remaining!")
						} else if tries == 3 {
							break
						}

						tries++
					}

					if password2 != admins[p].UserPassword {
						selector = "0"
						selector3 = "return"
					}

					selector = "3"
					selector3 = "0"

					break

				} else if !skipLogin {
					selector = "0"
					selector3 = "return"
				}
			}

			for selector3 != "1" && selector3 != "2" && selector3 != "3" && selector3 != "return" && selector3 != "end" {
				//password2 = ""
				fmt.Println(" ")
				fmt.Println(dataBase, " ")
				fmt.Println(" ")
				fmt.Printf("--------------------------\nTo edit existing entry, press 1.\n//\nTo remove an entry, press 2.\n//\nTo register admin/change password, press 3.\n//\nTo return, type \"return\".\n--------------------------\n")
				fmt.Scan(&selector3)

			}
			if selector3 == "end" {
				selector = "end"

			} else if selector3 == "return" {
				selector = "0"
				skipLogin = false

			} else if selector3 == "1" { // edit entry
				var idToChange int
				var propertyToChange, newValue string
				var updatePerson person

				fmt.Println("Enter Id number of entry to be modified:")
				fmt.Scan(&idToChange)

				for n := 0; n < len(dataBase); n++ {

					if dataBase[n].IdNumber == idToChange {

						updatePerson = dataBase[n]

						fmt.Println("Enter property to be modified:")
						//This doesnt work for multiple words e.g "favourite animal"
						//Fix this!
						fmt.Scan(&propertyToChange)

						fmt.Println("Enter new value:")
						fmt.Scan(&newValue)

						if propertyToChange == "name" {
							dataBase[n].Name = newValue
							updatePerson.Name = newValue

						} else if propertyToChange == "age" {
							newValueInt, err := strconv.Atoi(newValue)

							if err != nil {
								panic(err)
							}

							dataBase[n].Age = newValueInt
							updatePerson.Age = newValueInt

						} else if propertyToChange == "favourite animal" {
							dataBase[n].FavAnimal = newValue
							updatePerson.FavAnimal = newValue
						}
						_, err = col.Upsert(fmt.Sprintf("%d", updatePerson.IdNumber),
							person{
								Name:      updatePerson.Name,
								Age:       updatePerson.Age,
								FavAnimal: updatePerson.FavAnimal,
								IdNumber:  updatePerson.IdNumber,
							}, nil)
						if err != nil {
							log.Fatal(err)
						}
					}
				}

				selector = "3"
				selector3 = "0"
				skipLogin = true

			} else if selector3 == "2" {
				var idToDelete int
				fmt.Println("Enter Id number of person to be removed: ")
				fmt.Scan(&idToDelete)

				for l := 0; l < len(dataBase); l++ {

					if dataBase[l].IdNumber == idToDelete {
						dataBase[l] = dataBase[len(dataBase)-1]
						dataBase = dataBase[:len(dataBase)-1]
					}
				}

				cluster.Query(fmt.Sprintf("Delete from `%s`._default._default use KEYS $1 ", bucketName), &gocb.QueryOptions{PositionalParameters: []interface{}{fmt.Sprintf("%d", idToDelete)}})
				fmt.Printf("Entry %d has been removed\n \n", idToDelete)

				skipLogin = true

			} else if selector3 == "3" {
				var newAdmin loginCredentials
				var newUserName, newPassword string

				fmt.Println("Choose new username:")
				fmt.Scan(&newUserName)
				fmt.Println("Choose a password:")
				fmt.Scan(&newPassword)

				newAdmin.UserName = newUserName
				newAdmin.UserPassword = newPassword

				admins = append(admins, newAdmin)

				_, err = col2.Upsert(newAdmin.UserName,
					loginCredentials{
						UserName:     newAdmin.UserName,
						UserPassword: newAdmin.UserPassword,
					}, nil)
				if err != nil {
					log.Fatal(err)
				}

				selector = "3"
				selector3 = "0"
				skipLogin = true

			}

		}

	}
}
