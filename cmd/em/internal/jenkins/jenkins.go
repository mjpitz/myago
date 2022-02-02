// Copyright (C) 2022 Mya Pitzeruse
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package jenkins

import (
	"context"
	"net/http"
	"time"

	"github.com/mjpitz/myago/cmd/em/internal/index"
	"github.com/mjpitz/myago/zaputil"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type doc struct {
	Timestamp time.Time `gorm:"primaryKey"`
	Job       string    `gorm:"primaryKey"`
	Build     int       `gorm:"primaryKey;autoIncrement:false"`
	Stage     string    `gorm:"primaryKey"`
	Status    string
	Duration  int64
	Queued    int64
	Paused    int64
}

func (s doc) TableName() string {
	return "jenkins_stage"
}

type Config struct {
	BaseURL string           `json:"base_url" usage:"specify the base url of the jenkins instance we're indexing"`
	Jobs    *cli.StringSlice `json:"job"      usage:"provide an initial list of jobs to analyze"`
}

func Run(ctx context.Context, cfg Config, idx *index.Index) error {
	log := zaputil.Extract(ctx)

	log.Info("migrating db")
	err := idx.Migrate(doc{})
	if err != nil {
		return err
	}

	api := &API{
		baseURL: cfg.BaseURL,
		client:  http.DefaultClient,
	}

	jobs := cfg.Jobs.Value()
	if len(jobs) == 0 {
		log.Info("listing projects")

		listProjectsResp, err := api.Jobs().List(ctx)
		if err != nil {
			return err
		}

		for _, job := range listProjectsResp.Jobs {
			jobs = append(jobs, job.Name)
		}
	}

	docs := make([]interface{}, 0)

	for _, jobName := range jobs {
		log.Info("processing", zap.String("job", jobName))

		job, err := api.Jobs().Get(ctx, jobName)
		if err != nil {
			return err
		}

		for _, build := range job.Builds {
			pipeline, err := api.Pipelines().Get(ctx, jobName, build.Number)
			if err != nil {
				return err
			}

			if pipeline.Status == "IN_PROGRESS" {
				continue
			}

			for _, stage := range pipeline.Stages {
				if stage.Status == "NOT_EXECUTED" {
					continue
				}

				docs = append(docs, doc{
					Timestamp: time.UnixMilli(stage.StartTimeMillis),
					Job:       jobName,
					Build:     build.Number,
					Stage:     stage.Name,
					Status:    stage.Status,
					Duration:  stage.DurationMillis,
					Paused:    stage.PauseDurationMillis,
				})
			}

			docs = append(docs, doc{
				Timestamp: time.UnixMilli(pipeline.StartTimeMillis),
				Job:       jobName,
				Build:     build.Number,
				Stage:     "",
				Status:    pipeline.Status,
				Duration:  pipeline.DurationMillis,
				Queued:    pipeline.QueueDurationMillis,
				Paused:    pipeline.PauseDurationMillis,
			})
		}
	}

	log.Info("indexing", zap.Int("docs", len(docs)))
	return idx.Index(docs...)
}
