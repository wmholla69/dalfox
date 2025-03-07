package optimization

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"github.com/hahwul/dalfox/pkg/model"
)

// GenerateNewRequest is make http.Cilent
func GenerateNewRequest(url, payload string, options model.Options) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	if options.Data != "" {
		d := []byte(payload)
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer(d))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if options.Header != "" {
		h := strings.Split(options.Header, ": ")
		if len(h) > 1 {
			req.Header.Add(h[0], h[1])
		}
	}
	if options.Cookie != "" {
		req.Header.Add("Cookie", options.Cookie)
	}
	if options.UserAgent != "" {
		req.Header.Add("User-Agent", options.UserAgent)
	} else {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:75.0) Gecko/20100101 Firefox/75.0")
	}
	return req
}

// MakeHeaderQuery is generate http query with custom header
func MakeHeaderQuery(target, hn, hv string,options model.Options) (*http.Request, map[string]string) {
	tempMap := make(map[string]string)
	tempMap["type"] = "toBlind"
	tempMap["payload"] = hv
	tempMap["param"] = "thisisheadertesting"
	req, _ := http.NewRequest("GET", target, nil)
	if options.Data != "" {
		d := []byte("")
		req, _ = http.NewRequest("POST", target, bytes.NewBuffer(d))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if options.Header != "" {
		h := strings.Split(options.Header, ": ")
		if len(h) > 1 {
			req.Header.Add(h[0], h[1])
		}
	}
	if options.Cookie != "" {
		req.Header.Add("Cookie", options.Cookie)
	}
	if options.UserAgent != "" {
		req.Header.Add("User-Agent", options.UserAgent)
	} else {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:75.0) Gecko/20100101 Firefox/75.0")
	}

	req.Header.Add(hn, hv)
	return req, tempMap
}


// MakeRequestQuery is generate http query with custom paramters
func MakeRequestQuery(target, param, payload, ptype string, options model.Options) (*http.Request, map[string]string) {
	tempMap := make(map[string]string)
	tempMap["type"] = ptype
	tempMap["payload"] = payload
	tempMap["param"] = param

	payload = url.QueryEscape(payload)
	u, _ := url.Parse(target)
	data := u.String()
	if options.Data != "" {
		tempParam, _ := url.ParseQuery(options.Data)
		body := strings.Replace(options.Data, param+"="+tempParam[param][0], param+"="+tempParam[param][0]+payload, 1)
		tempURL, _ := url.Parse(data)
		tempQuery := tempURL.Query()
		tempURL.RawQuery = tempQuery.Encode()

		rst := GenerateNewRequest(tempURL.String(), body, options)
		return rst, tempMap

	} else {
		tempParam := u.Query()
		data = strings.Replace(data, param+"="+tempParam[param][0], param+"="+tempParam[param][0]+payload, 1)
		tempURL, _ := url.Parse(data)
		tempQuery := tempURL.Query()
		tempURL.RawQuery = tempQuery.Encode()
		rst := GenerateNewRequest(tempURL.String(), "", options)
		return rst, tempMap
	}
}

// MakePathQuery is generate http query with path
func MakePathQuery(target, fakeparam, payload, ptype string, options model.Options) (*http.Request, map[string]string) {
	tempMap := make(map[string]string)
	tempMap["type"] = ptype
	tempMap["payload"] = payload
	tempMap["param"] = fakeparam
	payload = url.QueryEscape(payload)
	u, err := url.Parse(target)
	if err != nil {
		rst := GenerateNewRequest(target, "", options)
		return rst, tempMap
	}
	data := ""
	if u.Path != "" {
		data = u.Scheme + "://" + u.Hostname() + u.Path + ";" + payload
	} else {
		data = u.Scheme + "://" + u.Hostname() + "/" + u.Path + ";" + payload
	}
	tempURL, err := url.Parse(data)
	if err != nil {
		rst := GenerateNewRequest(target, "", options)
		return rst, tempMap
	}

	tempQuery := tempURL.Query()
	tempURL.RawQuery = tempQuery.Encode()

	rst := GenerateNewRequest(tempURL.String(), options.Data, options)
	return rst, tempMap
}

// MakeURLEncodeRequestQuery is generate http query with Double URL Encoding
func MakeURLEncodeRequestQuery(target, param, payload, ptype string, options model.Options) (*http.Request, map[string]string) {

	tempMap := make(map[string]string)
	tempMap["type"] = ptype
	tempMap["payload"] = payload
	tempMap["param"] = param
	payload = url.QueryEscape(payload)

	u, _ := url.Parse(target)
	data := u.String()
	// URL Encoding
	encodedPayload := UrlEncode(UrlEncode(payload))
	if options.Data != "" {
		tempParam, _ := url.ParseQuery(options.Data)
		body := strings.Replace(options.Data, param+"="+tempParam[param][0], param+"="+tempParam[param][0]+encodedPayload, 1)
		tempURL, _ := url.Parse(data)
		tempQuery := tempURL.Query()
		tempURL.RawQuery = tempQuery.Encode()

		rst := GenerateNewRequest(tempURL.String(), body, options)
		return rst, tempMap

	} else {
		tempParam := u.Query()
		data = strings.Replace(data, param+"="+tempParam[param][0], param+"="+tempParam[param][0]+encodedPayload, 1)
		tempURL, _ := url.Parse(data)
		tempQuery := tempURL.Query()
		tempURL.RawQuery = tempQuery.Encode()
		rst := GenerateNewRequest(tempURL.String(), "", options)
		return rst, tempMap
	}
}

// MakeHTMLEncodeRequestQuery is generate http query with Hex Encoding
func MakeHTMLEncodeRequestQuery(target, param, payload, ptype string, options model.Options) (*http.Request, map[string]string) {
	tempMap := make(map[string]string)
	tempMap["type"] = ptype
	tempMap["payload"] = payload
	tempMap["param"] = param
	payload = url.QueryEscape(payload)

	u, _ := url.Parse(target)
	data := u.String()
	// HTML HEX Encoding
	encodedPayload := template.HTMLEscapeString(payload)
	if options.Data != "" {
		tempParam, _ := url.ParseQuery(options.Data)
		body := strings.Replace(options.Data, param+"="+tempParam[param][0], param+"="+tempParam[param][0]+encodedPayload, 1)
		tempURL, _ := url.Parse(data)
		tempQuery := tempURL.Query()
		tempURL.RawQuery = tempQuery.Encode()
		rst := GenerateNewRequest(tempURL.String(), body, options)
		return rst, tempMap

	} else {
		tempParam := u.Query()
		data = strings.Replace(data, param+"="+tempParam[param][0], param+"="+tempParam[param][0]+encodedPayload, 1)
		tempURL, _ := url.Parse(data)
		tempQuery := tempURL.Query()
		tempURL.RawQuery = tempQuery.Encode()
		rst := GenerateNewRequest(tempURL.String(), "", options)
		return rst, tempMap
	}
}

// Optimization is remove payload included badchar
func Optimization(payload string, badchars []string) bool {
	for _, v := range badchars {
		if strings.Contains(payload, v) {
			return false
		}
	}
	return true
}

// UrlEncode is custom url encoder for double url encoding
func UrlEncode(s string) (result string) {
	for _, c := range s {
		if c <= 0x7f { // single byte
			result += fmt.Sprintf("%%%X", c)
		} else if c > 0x1fffff { // quaternary byte
			result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
				0xf0+((c&0x1c0000)>>18),
				0x80+((c&0x3f000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else if c > 0x7ff { // triple byte
			result += fmt.Sprintf("%%%X%%%X%%%X",
				0xe0+((c&0xf000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else { // double byte
			result += fmt.Sprintf("%%%X%%%X",
				0xc0+((c&0x7c0)>>6),
				0x80+(c&0x3f),
			)
		}
	}

	return result
}
