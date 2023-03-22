package job

import (
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func NewJob(job interface{}, period *gocron.Job) {

	name := runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()

	scope := di.Scope(name)
	err := scope.Provide(job, dig.As(new(Job)))
	if err != nil {
		log.FatalWrap(err, "cannot initialize new job")
	}

	err = scope.Invoke(func(j Job) {
		err := period.Do(j.Run)
		if err != nil {
			log.Fatal("cannot DO cron %s", name)
		}
	})
	if err != nil {
		log.FatalWrap(err, "cannot create new job")
	}
	log.Info("New job was successfully registered - %s", name)
}

func Start() error {

	enabled := viper.GetBool("jobs.enabled")
	if !enabled {
		return nil
	}

	go gocron.Start()
	return nil
}

func Time(time string) *gocron.Job {

	split := strings.Split(time, " ")
	if len(split) != 2 {
		return nil
	}
	timeVal, err := strconv.ParseUint(split[0], 10, 64)
	if err != nil {
		return nil
	}
	timeId := split[1]

	job := gocron.Every(timeVal)
	switch timeId {
	case "seconds":
		job = job.Seconds()
	case "second":
		job = job.Second()
	case "minute":
		job = job.Minute()
	case "minutes":
		job = job.Minutes()
	case "hour":
		job = job.Hour()
	case "hours":
		job = job.Hours()
	}

	return job
}
