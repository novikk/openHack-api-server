package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"openhack/models"
	"os"
	"strconv"
)

var FIWARE_AUTH_TOKEN string

func init() {
	FIWARE_AUTH_TOKEN = os.Getenv("FIWARE_AUTH_TOKEN")
	fmt.Println(FIWARE_AUTH_TOKEN)
	GetPollutionData()
}

func GetPollutionData() {
	url := "http://orion.lab.fiware.org:1026/ngsi10/queryContext?limit=200"

	jsonStr := `{ "entities":[{ "type":"santander:device", "isPattern":"true", "id":"urn:x-iot:smartsantander:u7jcfa:mobile.*" }] }`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("X-Auth-Token", FIWARE_AUTH_TOKEN)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	type attr struct {
		Name  string
		Value string
	}

	type ctxElement struct {
		Attributes []attr `json:"attributes"`
	}

	type contextResponse struct {
		ContextElement ctxElement `json:"contextElement"`
	}

	type res struct {
		ContextResponses []contextResponse `json:"contextResponses"`
	}

	var r res
	json.Unmarshal(body, &r)

	var dt models.PollutionData

	db := models.GetDB()
	db.Exec("DELETE FROM pollution_points")

	for _, cr := range r.ContextResponses {
		var lat, lng, pollution float64
		ce := cr.ContextElement
		for _, at := range ce.Attributes {
			if at.Name == "Latitud" {
				lat, _ = strconv.ParseFloat(at.Value, 64)
			}

			if at.Name == "Longitud" {
				lng, _ = strconv.ParseFloat(at.Value, 64)
			}

			if at.Name == "NO2Concentration" {
				pollution, _ = strconv.ParseFloat(at.Value, 64)
			}
		}

		if pollution > 10 && pollution < 750 {
			p := models.PollutionPoint{
				Lat:   lat,
				Lng:   lng,
				Value: pollution,
			}
			db.Create(&p)
			dt = append(dt, p)
		}
	}

	type finalRes struct {
		Points models.PollutionData `json:"points"`
	}

	fr := finalRes{
		Points: dt,
	}

	rrr, _ := json.Marshal(fr)
	//fmt.Println(string(rrr[:]))
	url = fmt.Sprintf("http://localhost:8080/api/newPollution?pol=%s", string(rrr[:]))

	http.Get(url)
}
