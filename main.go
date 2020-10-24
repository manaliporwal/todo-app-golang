package main

func main() {
	initDB()
	initSessionCache()
	//session, _ := addSession("manali", "23/56/2020")
	//fmt.Printf(session)
	//result := validateSession("manali", "123")
	//if (result == true) {
	//	fmt.Printf("True")
	//}
	initRouter()
}
