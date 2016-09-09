package tvdbapi

import(
	"fmt"
	"net/http"
	"bytes"
	"log"
	"io/ioutil"
	"encoding/json"
)

type TvDbConfig struct {
	Username string
	ApiKey string
	UserKey string
}

var apitoken = ""

func Login(config TvDbConfig)  {
	client := http.Client{}

	jsonStr := []byte(fmt.Sprintf(`{"username": "%s", "apikey": "%s", "userkey": "%s"}`, config.Username, config.ApiKey, config.UserKey))
	resp, err := client.Post("https://api.thetvdb.com/login", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var m interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}

	if apitokentmp, ok := m.(map[string]interface{})["token"]; ok {
		apitoken = apitokentmp.(string)
		fmt.Println("login success")
	} else {
		fmt.Println("login unsuccess")
	}
}