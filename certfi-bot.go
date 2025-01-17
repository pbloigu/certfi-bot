package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/form/v4"
	"github.com/google/uuid"
	"github.com/pbloigu/gonfig"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type configStruct struct {
	RunAs struct {
		Uid int `yaml:"uid"`
		Gid int `yaml:"gid"`
	} `yaml:"runAs"`
	Schedule struct {
		IntervalHours int `yaml:"intervalHours"`
	} `yaml:"schedule"`
	Feed struct {
		Url string `yaml:"url"`
	} `yaml:"feed"`
	Server struct {
		Host        string `yaml:"url"`
		AccessToken string `yaml:"accessToken"`
	} `yaml:"server"`
}

func (c *configStruct) Uid() int {
	return c.RunAs.Uid
}

func (c *configStruct) Gid() int {
	return c.RunAs.Gid
}

func (c *configStruct) StatePath() string {
	return ""
}

type toot struct {
	Status     string `form:"status"`
	Visibility string `form:"visibility"`
	Language   string `form:"language"`
}

var config configStruct

var encoder = form.NewEncoder()

func main() {
	setupLoggig()
	if err := gonfig.Get("/etc/certfi-bot/config.yml", &config); err != nil {
		log.Panic().Err(err).Msg("Startup failed.")
	}

	s := getScheduler()

	createJob(s)

	s.Start()
	waitForTermination()

	err := s.Shutdown()
	log.Info().Msg("Scheduler shut down.")
	if err != nil {
		log.Error().AnErr("error", err).Msg("Scheduler stop failed.")
	}
}

func waitForTermination() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func setupLoggig() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled.")
	}
}

func getScheduler() gocron.Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Scheduler init failed.")
		os.Exit(-1)
	}
	log.Info().Msg("Scheduler initialized.")
	return s
}

func createJob(s gocron.Scheduler) gocron.Job {

	interval := time.Duration(config.Schedule.IntervalHours) * time.Hour
	log.Info().Any("interval", interval).Msg("Execution interval set.")
	job, err := s.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(runner, interval),
	)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("Job init failed.")
		os.Exit(-1)
	}
	log.Info().Str("jobId", job.ID().String()).Msg("Job created.")
	return job
}

func runner(interval time.Duration) {
	log.Info().Msg("Fetching feed.")
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(config.Feed.Url)
	if err != nil {
		log.Error().AnErr("error", err).Msg("Parsing feed failed. Skipping.")
		return
	}
	log.Debug().Str("feedTitle", feed.Title).Msg("Parsed feed.")
	for _, item := range feed.Items {
		pubTime, err := time.Parse(time.RFC1123, item.Published)
		if err != nil {
			log.Error().AnErr("error", err).Str("itemTitle", item.Title).Msg("Parsing publish time failed. Skipping this feed item.")
			continue
		} else {
			if isNew(pubTime, interval) {
				toot := createToot(*item)
				log.Debug().Any("toot", toot).Msg("Created a toot.")
				request := createRequest(toot)
				doToot(request)
			}
		}
	}
}

func isNew(pubTime time.Time, interval time.Duration) bool {
	return time.Since(pubTime) < interval
}

func createToot(item gofeed.Item) toot {
	return toot{
		Status:     fmt.Sprintf("%s\n%s\n\n%s", item.Title, item.Description, item.Link),
		Language:   "en",
		Visibility: "public",
	}
}

func doToot(r http.Request) {
	client := &http.Client{}
	client.Do(&r)
}

func createRequest(t toot) http.Request {
	body, _ := encoder.Encode(&t)
	bodyReader := bytes.NewReader([]byte(body.Encode()))
	request, _ := http.NewRequest("POST", config.Server.Host+"/api/v1/statuses", bodyReader)
	request.Header.Add("Authorization", "Bearer "+config.Server.AccessToken)
	request.Header.Add("Idempotency-Key", uuid.NewString())
	return *request
}
