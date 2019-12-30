package main

import "net/http"

import "io/ioutil"

import "time"

import "log"

type Stats struct {
	Success  int    `json:"success,omitempty"`
	Errors   int    `json:"errors,omitempty"`
	LastBody string `json:"last_body,omitempty"`
}

type Client struct {
	client *http.Client
	target string
	rate   int

	stats Stats
	Stats chan Stats

	throttle <-chan time.Time
}

func NewClient(target string, rate int) *Client {
	return &Client{
		client: &http.Client{
			//Timeout: time.Second,
		},
		target:   target,
		Stats:    make(chan Stats),
		throttle: time.Tick(time.Second / time.Duration(rate)),
	}
}

func (r *Client) Run() {
	for i := 0; i < 200; i++ {
		go r.loop()
	}

	select {}
}

func (r *Client) loop() {
	for range r.throttle {
		r.req()
	}
}

func (r *Client) req() {
	defer r.sendStats()
	start := time.Now()
	defer func() {
		log.Printf("requested %s (%s)", r.target, time.Since(start))
	}()

	req, _ := http.NewRequest(http.MethodGet, r.target, nil)

	res, err := r.client.Do(req)
	if err != nil {
		log.Printf("error requesting %s: %s", r.target, err)
		r.stats.LastBody = "Error"
		r.stats.Errors++
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		log.Printf("error requesting %s: code %d", r.target, res.StatusCode)
		r.stats.LastBody = "Error"
		r.stats.Errors++
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		r.stats.LastBody = "Error"
		r.stats.Errors++
		return
	}

	r.stats.Success++
	r.stats.LastBody = string(body)
}

func (r *Client) sendStats() {
	select {
	case r.Stats <- r.stats:
	default:
	}
}
