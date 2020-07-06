package bw

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)



const (
	statsPath                   = "/wp-admin/admin-ajax.php"
	contentTypeWWWormUrlEncoded = "application/x-www-form-urlencoded"
)

var (
	statsBody *url.Values
)

func init() {
	statsBody = &url.Values{}
	statsBody.Set("action", "cxo_get_crowd_indicator")
}

func StatsForBW(domain string) (stats Stats, err error) {
	var resp *http.Response
	if resp, err = http.Post(
		fmt.Sprintf("https://%s%s", domain, statsPath),
		contentTypeWWWormUrlEncoded,
		strings.NewReader(statsBody.Encode()),
	); err != nil {
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&stats)
	return
}
