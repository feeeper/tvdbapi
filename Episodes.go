package tvdbapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type episodeData struct {
	Episodes []Episode `json:"data"`
}

type Episode struct {
	AbsoluteNumber     int       `json:"absoluteNumber"`
	AiredEpisodeNumber int       `json:"airedEpisodeNumber"`
	AiredSeason        int       `json:"airedSeason"`
	DvdEpisodeNumber   float32   `json:"dvdEpisodeNumber"`
	DvdSeason          int       `json:"dvdSeason"`
	EpisodeName        string    `json:"episodeName"`
	Id                 int       `json:"id"`
	Overview           string    `json:"overview"`
	FirstAired         AiredTime `json:"firstAired"`
}

func (client Client) GetEpisodes(series Series) []Episode {
	return client.GetEpisodesBySeriesId(series.Id)
}

func (client Client) GetEpisodesBySeriesId(seriesId int) []Episode {
	result := episodeData{}

	url := fmt.Sprintf("https://api.thetvdb.com/series/%v/episodes", seriesId)

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