package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Netflix struct{}

type NetflixPostingsInfo struct {
	Pages int `json:"num_pages"`
}

type NetflixInfo struct {
	PostingsInfo NetflixPostingsInfo `json:"postings"`
}

type NetflixPosting struct {
	Title    string `json:"text"`
	Desc     string `json:"description"`
	Location string `json:"location"`
	Url      string `json:"url"`
}

type NetflixRecords struct {
	Postings []NetflixPosting `json:"postings"`
}

type NetflixSearch struct {
	Info    NetflixInfo    `json:"info"`
	Records NetflixRecords `json:"records"`
}

func (netflix *Netflix) readPage(url string) ([]*Job, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s\n", res.StatusCode, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	search := NetflixSearch{}
	err = json.Unmarshal(body, &search)
	if err != nil {
		log.Fatal(err)
	}

	jobs := make([]*Job, 0)

	for _, p := range search.Records.Postings {
		jobs = append(jobs, &Job{
			Title: p.Title,
			Location: p.Location,
			Type: string(FullTime),
			Desc: p.Desc,
			Link: p.Url,
		})
	}

	return jobs, nil
}

func (netflix *Netflix) ListJobs() []*Job {
	jobs := make([]*Job, 0)

	res, err := http.Get("https://jobs.netflix.com/api/search")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s\n", res.StatusCode, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	search := NetflixSearch{}
	err = json.Unmarshal(body, &search)
	if err != nil {
		log.Fatal(err)
	}

	for i:=1; i<=search.Info.PostingsInfo.Pages; i++ {
		j, err := netflix.readPage(fmt.Sprintf("https://jobs.netflix.com/api/search?page=%d", i))
		if err != nil {
			log.Fatal(err)
		}

		jobs = append(jobs, j...)
	}

	return jobs
}