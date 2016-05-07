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

	var r interface{}
	json.Unmarshal(body, &r)

	c.Data["json"] = r
}
