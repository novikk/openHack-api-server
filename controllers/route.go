package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

type RouteController struct {
	beego.Controller
}

func (c *RouteController) NewRoute() {
	defer c.ServeJSON()

	tp := c.GetString("tp")
	km := c.GetString("km")
	lat := c.GetString("fromLat")
	lon := c.GetString("fromLon")

	url := fmt.Sprintf("http://localhost:8080/api/newRoute?tp=%s&km=%s&fromLat=%s&fromLon=%s", tp, km, lat, lon)
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		c.Data["json"] = "error"
		return
	}

	r, _ := ioutil.ReadAll(res.Body)
	var rr interface{}
	json.Unmarshal(r, &rr)
	c.Data["json"] = rr
}
