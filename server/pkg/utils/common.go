package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Base64Decode(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	return string(b), err
}

func DoGetRequest(urlStr, token string, params map[string]any) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		if v != nil {
			q.Add(k, fmt.Sprint(v))
		}
	}
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 40 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return string(body), nil
}

func DoPostRequest(urlStr, token string, payload any) (string, error) {
	b, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 40 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)
	return string(respBody), nil
}

type Message struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Response struct {
	Data any     `json:"data,omitempty"`
	Mess Message `json:"mess"`
}

func ResponseSuccess(data any) Response {
	return Response{
		Data: data,
		Mess: Message{Code: 200, Description: "SUCCESS"},
	}
}

func ResponseError(code int, desc string) Response {
	return Response{
		Mess: Message{Code: code, Description: desc},
	}
}

func DownloadFile(c *gin.Context, filePath, fileName string) {
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

func RandomPassword() string {
	lower := "qwertyuiopasdfghjklzxcvbnm"
	upper := strings.ToUpper(lower)
	num := "0123456789"
	special := "!@#$"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	gen := []byte{
		lower[r.Intn(len(lower))],
		upper[r.Intn(len(upper))],
		num[r.Intn(len(num))],
		special[r.Intn(len(special))],
	}

	all := lower + upper + num + special
	for len(gen) < 8 {
		gen = append(gen, all[r.Intn(len(all))])
	}

	r.Shuffle(len(gen), func(i, j int) {
		gen[i], gen[j] = gen[j], gen[i]
	})

	return string(gen)
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
