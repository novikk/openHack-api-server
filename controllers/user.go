package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"openhack/models"
	"strconv"
)

type UserController struct {
	beego.Controller
}

func newCustomer(uid string) string {
	url := "http://192.176.47.48:27030/rest/S-HcN8-IUGtV-/customerManagement/v2/customer"

	var jsonStr = []byte(`
		{
		  "name": "` + uid + `",
		  "status": "Active",
		  "characteristic": [
		    {
		      "name": "points",
		      "value": "0"
		    }
		  ]
		}
	`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", TMF_CUSTOMER_KEY))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	type response struct {
		ID int `json:"id"`
	}

	var r response
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &r)

	db := models.GetDB()
	db.Create(&models.User{
		CID:  r.ID,
		UUID: uid,
	})

	return string(body[:])
}

func findCustomer(id int) string {
	url := fmt.Sprintf("http://192.176.47.48:27030/rest/S-HcN8-IUGtV-/customerManagement/v2/customer/%d", id)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", TMF_CUSTOMER_KEY))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body[:])
}

func getPoints(js string) int {
	fmt.Println(js)
	type charact struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	type respts struct {
		Characteristics []charact `json:"characteristic"`
	}

	var r respts
	json.Unmarshal([]byte(js), &r)

	fmt.Println(r)

	pts, _ := strconv.Atoi(r.Characteristics[0].Value)
	return pts
}

func setPoints(id, pts int) {
	url := fmt.Sprintf("http://192.176.47.48:27030/rest/S-HcN8-IUGtV-/customerManagement/v2/customer/%d", id)

	var jsonStr = []byte(`
		{
		  "characteristic": [
		    {
		      "name": "points",
		      "value": "` + strconv.Itoa(pts) + `"
		    }
		  ]
		}
	`)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", TMF_CUSTOMER_KEY))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}

func (c *UserController) FindOrCreateUser() {
	db := models.GetDB()
	uid := c.GetString("uid")

	user := models.User{
		UUID: "",
	}
	db.First(&user, "UUID = ?", uid)

	if user.UUID != "" {
		c.Data["json"] = findCustomer(user.CID)
	} else {
		c.Data["json"] = newCustomer(uid)
	}

	var r interface{}
	json.Unmarshal([]byte(c.Data["json"].(string)), &r)
	c.Data["json"] = r
	c.ServeJSON()
}

func (c *UserController) AddPoints() {
	db := models.GetDB()
	uid := c.GetString("uid")
	newpts, _ := c.GetInt("pts")

	user := models.User{
		UUID: "",
	}
	db.First(&user, "UUID = ?", uid)

	cust := findCustomer(user.CID)
	pts := getPoints(cust)
	pts += newpts

	setPoints(user.CID, pts)
	c.Data["json"] = "ok"
	c.ServeJSON()
}

func (c *UserController) Buy() {
	defer c.ServeJSON()
	db := models.GetDB()
	uid := c.GetString("uid")
	pOffering := c.GetString("offering")
	price, _ := c.GetInt("price") // TODO: get price from product on server

	user := models.User{
		UUID: "",
	}
	db.First(&user, "UUID = ?", uid)

	udata := findCustomer(user.CID)
	cpts := getPoints(udata)

	if cpts < price {
		c.Data["json"] = "Not enough points"
		return
	}

	cpts -= price
	setPoints(user.CID, cpts)

	url := "http://192.176.47.48:27030/rest/S-LcN8-IUGtV-/productInventory/v2/product"

	var jsonStr = []byte(`
		{
	    "name": "Item bought",
	    "description": "Bought product",
	    "status": "Created",
	    "isCustomerVisible": true,
	    "isBundle" : true,
	    "productSerialNumber": "GEN_QR",
	    "startDate": "2014-04-25T12:16:43.0Z",
	    "orderDate": "2014-04-25T12:16:43.0Z",
	    "terminationDate": "",
	    "productOffering":
	    {
        "id": "http://192.176.47.48:27030/catalogApi/productOffering/` + pOffering + `",
        "name": "` + pOffering + `"
	    },

	    "relatedParty": [
	    {
        "id": "` + strconv.Itoa(user.CID) + `",
        "role":"owner",
        "href":""
	    }]
	}
	`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", TMF_INVENTORY_KEY))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp)

	c.Data["json"] = "ok"
}

func (c *UserController) GetCustomerInventory() {
	defer c.ServeJSON()
	db := models.GetDB()
	uid := c.GetString("uid")

	user := models.User{
		UUID: "",
	}
	db.First(&user, "UUID = ?", uid)

	type prodoff struct {
		Name string `json:"name"`
	}

	type party struct {
		ID string `json:"id"`
	}

	type prodresp struct {
		QR       string  `json:"productSerialNumber"`
		Offering prodoff `json:"productOffering"`
		Parties  []party `json:"relatedParty"`
	}

	url := "http://192.176.47.48:27030/rest/S-LcN8-IUGtV-/productInventory/v2//product"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", TMF_INVENTORY_KEY))
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

	var pr []prodresp
	json.Unmarshal(body, &pr)

	type result struct {
		Name string
		QR   string
	}

	var res []result
	for _, p := range pr {
		if p.Parties[0].ID == strconv.Itoa(user.CID) {
			res = append(res, result{
				Name: p.Offering.Name,
				QR:   p.QR,
			})
		}
	}

	c.Data["json"] = res
}
