package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type NameResponse struct {
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value struct {
		Id   int    `json:"id"`
		Joke string `json:"joke"`
	} `json:"value"`
}

func main() {
	// command to install: go get github.com/gin-gonic/gin
	//initialises a router with the default functions.
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		//context.String(http.StatusOK, "Hello world!")
		var name = new(NameResponse)
		callservice("https://names.mcquay.me/api/v0/", name)

		//parse the url.
		urlA, err := url.Parse("http://api.icndb.com/jokes/random")
		if err != nil {
			fmt.Print(err.Error())
		}

		//Adding parameters to url.
		values := urlA.Query()
		values.Add("firstName", name.First_Name)
		values.Add("lastName", name.Last_Name)
		values.Add("limitTo", "nerdy")

		urlA.RawQuery = values.Encode()
		fmt.Println(urlA.String())

		var joke = new(JokeResponse)
		callservice(urlA.String(), joke)

		//joke.Value.Joke = strings.Replace(joke.Value.Joke, "John Doe", fmt.Sprintf("%s %s", name.First_Name, name.Last_Name), 1)

		//context.IndentedJSON(http.StatusOK, joke)
		context.String(http.StatusOK, joke.Value.Joke)
	})

	// starts the server at port 8080
	router.Run("localhost:8080")
}

func getservicebody(url string) []byte {
	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	return bodyBytes
}

func callservice(url string, model interface{}) {
	bodyBytes := getservicebody(url)
	json.Unmarshal(bodyBytes, &model)
	fmt.Printf("API Response as struct %+v\n", model)
}
