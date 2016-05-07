package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

type ProductController struct {
	beego.Controller
}

func (c *ProductController) GetAllOfferings() {
	defer c.ServeJSON()
	url := fmt.Sprintf("http://192.176.47.48:27030/rest/CatalogManagement/v2/productOffering/?project=%s", beego.AppConfig.String("tmf_project"))

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", TMF_CATALOG_KEY))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		c.Data["json"] = ""
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body[:]))

	type res struct {
		ID string `json:"id"`
	}

	type price struct {
		Amount int `json:"amount"`
	}

	type ver struct {
		Price price `json:"price"`
	}

	type VP struct {
		Versions []ver `json:"versions"`
	}

	var rr []res
	json.Unmarshal(body, &rr)

	type product struct {
		Name  string
		Price int
	}

	var products []product

	for _, r := range rr {
		url := fmt.Sprintf("http://192.176.47.48:27030/rest/CatalogManagement/v2/productOffering/%s/productOfferingPrice/ctOneTime/?project=%s", r.ID, beego.AppConfig.String("tmf_project"))

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", TMF_CATALOG_KEY))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			c.Data["json"] = ""
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		var rsp []VP
		err = json.Unmarshal(body, &rsp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		products = append(products, product{
			Name:  r.ID,
			Price: rsp[0].Versions[0].Price.Amount,
		})
	}

	c.Data["json"] = products
}
