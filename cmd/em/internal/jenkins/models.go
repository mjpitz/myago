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

type ListJobsResponse struct {
	 Jobs []Job `json:"jobs"`
}

type Job struct {
	Name string `json:"name"`
}

type GetJobResponse struct {
	Builds []Build `json:"builds"`
}

type Build struct {
	Number int `json:"number"`
}

type GetPipelineResponse struct {
	ID                  string          `json:"id"`
	Status              string          `json:"status"`
	StartTimeMillis     int64           `json:"startTimeMillis"`
	DurationMillis      int64           `json:"durationMillis"`
	QueueDurationMillis int64           `json:"queueDurationMillis"`
	PauseDurationMillis int64           `json:"pauseDurationMillis"`
	Stages              []PipelineStage `json:"stages"`
}

type PipelineStage struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Status              string `json:"status"`
	StartTimeMillis     int64  `json:"startTimeMillis"`
	DurationMillis      int64  `json:"durationMillis"`
	PauseDurationMillis int64  `json:"pauseDurationMillis"`
}
