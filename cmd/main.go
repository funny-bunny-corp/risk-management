package main

func main() {
	mngr, err := buildAppContainer()
	if err != nil {
		panic(err)
	}
	mngr.Start()
}
