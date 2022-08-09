package main

func main() {
	a := App{}
	a.Initialize("root","sss","go-gorilla-rest-api-swagger")
	a.Run(":8090")
}