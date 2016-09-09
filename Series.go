package tvdbapi

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/url"
)

type searchData struct {
	Series []Series  `json:"data"`
}

type SearchQuery struct {
	Name string
	ImdbId string
	Zap2itId string
	AcceptLanguage string
}

type Series struct {
	Id int `json:"id"`
	SeriesName string `json:"seriesName"`
	Aliases []string `json:"aliases"`
	Banner string `json:"banner"`
	SeriesId string `json:"seriesId"`
	Status string `json:"status"`
	FirstAired string `json:"firstAired"`
	Network string `json:"network"`
	NetworkId string `json:"networkId"`
	Runtime string `json:"runtime"`
	Genre []string `json:"genre"`
	Overview string `json:"overview"`
	LastUpdated int `json:"lastUpdated"`
	AirsDayOfWeek string `json:"airsDayOfWeek"`
	AirsTime string `json:"airsTime"`
	Rating string `json:"rating"`
	ImdbId string `json:"imdbId"`
	Zap2itId string `json:"zap2itId"`
	Added string `json:"added"`
	SiteRating float32 `json:"siteRating"`
	SiteRatingCount int `json:"siteRatingCount"`
}

type seriesInfoData struct {
	SeriesInfo Series `json:"data"`
}

func Search(query SearchQuery) ([]Series) {
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

	req.Header.Add("authorization", "Bearer " + apitoken)
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

func GetSeriesInfo(series Series) Series {
	return GetSeriesInfoById(series.Id)
}

func GetSeriesInfoById(seriesId int) Series {
	result := seriesInfoData{}

	url := fmt.Sprintf("https://api.thetvdb.com/series/%v", seriesId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer " + apitoken)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	err := json.Unmarshal(body, &result)

	if err != nil {
		log.Fatal(err)
		return result.SeriesInfo
	}

	log.Println("get series info completed successfully")
	log.Println(fmt.Sprintf("Series: %s; ImdbId: %s; LastUpdated: %s; Zip2itid: %s",
		result.SeriesInfo.SeriesName,
		result.SeriesInfo.ImdbId,
		result.SeriesInfo.LastUpdated,
		result.SeriesInfo.Zap2itId))

	return result.SeriesInfo
}