package tvdbapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type searchData struct {
	Series []Series `json:"data"`
}

type SearchQuery struct {
	Name           string
	ImdbId         string
	Zap2itId       string
	AcceptLanguage string
}

type AiredTime struct {
	time.Time
}

const ctLayout = "2006-01-02"

func (ct *AiredTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	if len(b) == 0 {
		// sometimes AiredDate is empty. Lets it will be 1900-01-01
		ct.Time, err = time.Parse(ctLayout, "2999-12-31")
	} else {
		ct.Time, err = time.Parse(ctLayout, string(b))
	}
	return
}

type Series struct {
	Id              int       `json:"id"`
	SeriesName      string    `json:"seriesName"`
	Aliases         []string  `json:"aliases"`
	Banner          string    `json:"banner"`
	SeriesId        string    `json:"seriesId"`
	Status          string    `json:"status"`
	FirstAired      AiredTime `json:"firstAired"`
	Network         string    `json:"network"`
	NetworkId       string    `json:"networkId"`
	Runtime         string    `json:"runtime"`
	Genre           []string  `json:"genre"`
	Overview        string    `json:"overview"`
	LastUpdated     int       `json:"lastUpdated"`
	AirsDayOfWeek   string    `json:"airsDayOfWeek"`
	AirsTime        string    `json:"airsTime"`
	Rating          string    `json:"rating"`
	ImdbId          string    `json:"imdbId"`
	Zap2itId        string    `json:"zap2itId"`
	Added           string    `json:"added"`
	SiteRating      float32   `json:"siteRating"`
	SiteRatingCount int       `json:"siteRatingCount"`
}

type seriesInfoData struct {
	Series Series `json:"data"`
}

type Update struct {
	Updated		time.Time	`json:"lastUpdated"`
	SeriesId	int		`json:""id`
}

type updates struct {
	Updates		[]Update	`json:"data"`
}

func (client Client) Search(query SearchQuery) []Series {
	result := searchData{}
	values := url.Values{}

	if query.Name != "" {
		values.Add("name", query.Name)
	}

	if query.ImdbId != "" {
		values.Add("imdbId", query.ImdbId)
	}

	if query.Zap2itId != "" {
		values.Add("zap2itId", query.Zap2itId)
	}

	url := fmt.Sprintf("https://api.thetvdb.com/search/series?%s", values.Encode())

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+client.ApiToken)
	if query.AcceptLanguage != "" {
		req.Header.Add("Accept-Language", query.AcceptLanguage)
	}

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	err := json.Unmarshal(body, &result)

	if err != nil {
		log.Fatal(err)
		return result.Series
	}

	log.Println("search completed successfully")
	log.Println(fmt.Sprintf("Total found: %v", len(result.Series)))

	return result.Series
}

func (client Client) GetSeriesInfo(series Series) Series {
	return client.GetSeriesInfoById(series.Id)
}

func (client Client) GetSeriesInfoById(seriesId int) Series {
	result := seriesInfoData{}

	url := fmt.Sprintf("https://api.thetvdb.com/series/%v", seriesId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+client.ApiToken)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	err := json.Unmarshal(body, &result)

	if err != nil {
		log.Fatal(err)
		return result.Series
	}

	log.Println("get series info completed successfully")
	log.Println(fmt.Sprintf("Series: %s; ImdbId: %s; LastUpdated: %s; Zip2itid: %s; First Aired: %v",
		result.Series.SeriesName,
		result.Series.ImdbId,
		result.Series.LastUpdated,
		result.Series.Zap2itId,
		result.Series.FirstAired))

	return result.Series
}

func (client Client) GetUpdate(fromTime time.Time) ([]Update, error) {
	var result []Update
	var updates updates
	url := fmt.Sprintf("https://api.thetvdb.com/updated/query?fromTime=%s", fromTime.Format("UnixTime"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("authorization", "Bearer " + client.ApiToken)

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return result, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &updates)
	if err != nil {
		return result, err
	}

	return updates.Updates, nil
}