package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Define struct types for unmarshaling API response
type Item struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	CreatedDate string `json:"createdDate"`
	Url         string `json:"url"`
	Body        string `json:"body"`
	MainArtist  Artist `json:"mainArtist"`
	Tags        []Tag  `json:"tags"`
}

type Tag struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

type Artist struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

type ApiResponse struct {
	Content []Item `json:"content"`
}

func main() {
	apiURL := "https://akordi.lv/api/songs?includeBody=true&size=5000" // Replace with your JSON API endpoint
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error fetching data from API:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading API response:", err)
		return
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, item := range apiResponse.Content {
		filename := "content/song/" + getLastPartOfURL(item.Url) + ".md"
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating file '%s': %v\n", filename, err)
			continue
		}

		formattedDate := convertDate(item.CreatedDate)

		var tags []string
		for _, tag := range item.Tags {
			tags = append(tags, fmt.Sprintf(`"%s"`, escapeQuotes(tag.Title)))
		}
		tagsJoined := strings.Join(tags, ", ")

		content := fmt.Sprintf("---\n"+
			"title: \"%s\"\n"+
			"date: %s\n"+
			"url: %s\n"+
			"categories: [\"%s\"]\n"+
			"tags: [%s]\n"+
			"draft: false\n"+
			"---\n"+
			"```text\n%s\n```",
			escapeQuotes(item.Title),
			formattedDate,
			item.Url,
			escapeQuotes(ReplaceDots(item.MainArtist.Title)),
			tagsJoined,
			item.Body)

		_, err = file.WriteString(content)
		file.Close()

		if err != nil {
			fmt.Printf("Error writing to file '%s': %v\n", filename, err)
		} else {
			fmt.Printf("Created file '%s' successfully\n", filename)
		}
	}
}

func getLastPartOfURL(url string) string {
	parts := strings.Split(url, "/")
	lastPart := parts[len(parts)-1]
	return lastPart
}

func convertDate(dateStr string) string {
	layoutISO := "2006-01-02T15:04:05.000-07:00"
	layoutRFC := "2006-01-02T15:04:05-07:00"
	t, err := time.Parse(layoutISO, dateStr)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return t.Format(layoutRFC)
}

func escapeQuotes(str string) string {
	escapedStr := strings.ReplaceAll(str, "\"", "\\\"")
	return escapedStr
}

func ReplaceDots(str string) string {
	return strings.ReplaceAll(str, ".", "-")
}
