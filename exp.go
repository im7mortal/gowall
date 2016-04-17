package main

import "time"

/*import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	session, err := mgo.Dial("mongodb://localhost:27017/lol")
	defer session.Close()
	if err != nil {
		println(err.Error())
	}
	d := session.DB("")
	collection := d.C("hj")
	i, err := collection.Count()

	println(i)
	_ = collection.Insert(bson.M{"hj":"ki"})
	if err != nil {
		println(err.Error())
	}
}*/
var fir = make(chan bool)
var sec = make(chan bool)
func sd() {
	time.Sleep(time.Duration(time.Second * 4))
	println("2")
	fir <- true
}
func sds() {
	time.Sleep(time.Duration(time.Second * 8))
	println("1")
	sec <- true
}
func main()  {


	go  sd()
	go sds()
	gh, hj := <-fir, <- sec
	println("dsdssdds")
	println(gh, hj)
}
