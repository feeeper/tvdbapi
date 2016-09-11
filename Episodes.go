package tvdbapi

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
)

type episodeData struct {
	Episodes []Episode `json:"data"`
}

type Episode struct {
	AbsoluteNumber int `json:"absoluteNumber"`
	AiredEpisodeNumber int `json:"airedEpisodeNumber"`
	AiredSeason int `json:"airedSeason"`
	DvdEpisodeNumber int `json:"dvdEpisodeNumber"`
	DvdSeason int `json:"dvdSeason"`
	EpisodeName string `json:"episodeName"`
	Id int `json:"id"`
	Overview string `json:"overview"`
	FirstAired string `json:"firstAired"`
}

func (client Client) GetEpisodes(series Series) []Episode {
	result := episodeData{}

	url := fmt.Sprintf("https://api.thetvdb.com/series/%v/episodes", series.Id)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer " + client.ApiToken)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	err := json.Unmarshal(body, &result)

	if err != nil {
		log.Fatal(err)
		return result.Episodes
	}

	log.Println("get episodes completed successfully")
	log.Println(fmt.Sprintf("Total episodes: %v", len(result.Episodes)))

	return result.Episodes
}