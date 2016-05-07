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
