package main

func main() {
	a := App{}
	a.Initialize("root","","127.0.0.1","3306","go-gorilla-rest-api-swagger")
	a.Run(":8090")
}