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

type Client struct {
	ApiToken string
}

func GetClient(config TvDbConfig) (Client, error) {
	result := Client{}
	client := http.Client{}

	jsonStr := []byte(fmt.Sprintf(`{"username": "%s", "apikey": "%s", "userkey": "%s"}`, config.Username, config.ApiKey, config.UserKey))
	resp, err := client.Post("https://api.thetvdb.com/login", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return result, err
	}

	var m interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
		return result, err
	}

	if apiToken, ok := m.(map[string]interface{})["token"]; ok {
		result.ApiToken = apiToken.(string)
		fmt.Println("login success")
		return result, nil
	} else {
		fmt.Println("login unsuccess")
		return result, LoginFailure{}
	}
}

type LoginFailure struct {}
func (lf LoginFailure) Error() string {
	return "login failure"
}