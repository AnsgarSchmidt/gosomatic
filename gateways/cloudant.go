package gateways

import (
	"github.com/smithsz/go-cloudant"
	"fmt"
	"os"
)

type Doc struct {
    Foo    string    `json:"foo"`
}

func Cloudant(){

	client, err := cloudant.CreateClient("3a4d4cf0-aed2-4916-8413-fa0177d2129f-bluemix", "24b4187fbd39510e84cc2cf10184cebf97ea56b836aab8ce4590ffe6477ae925", "https://3a4d4cf0-aed2-4916-8413-fa0177d2129f-bluemix.cloudant.com", 5)
	if err != nil{
		fmt.Println(err)
		os.Exit(2)
	}

	db, err := client.GetOrCreate("godb")
	if err != nil{
		fmt.Println(err)
		os.Exit(2)
	}

	myDoc := &Doc{
		Foo:    "bar",
	}

	newRev, err := db.Set(myDoc)
	if err != nil{
		fmt.Println(err)
		os.Exit(3)
	}

	fmt.Println(newRev)  // prints '_rev' of new document revision

}