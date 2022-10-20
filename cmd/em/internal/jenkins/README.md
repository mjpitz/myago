# jenkins




```go
import go.pitz.tech/lib/cmd/em/internal/jenkins
```

## Usage

#### func  Run

```go
func Run(ctx context.Context, cfg Config, idx *index.Index) error
```

#### type API

```go
type API struct {
}
```


#### func (*API) Do

```go
func (api *API) Do(ctx context.Context, method, url string, body io.Reader) (io.ReadCloser, error)
```

#### func (*API) Jobs

```go
func (api *API) Jobs() *Jobs
```

#### func (*API) Pipelines

```go
func (api *API) Pipelines() *Pipelines
```

#### type Build

```go
type Build struct {
	Number int `json:"number"`
}
```


#### type Config

```go
type Config struct {
	BaseURL string           `json:"base_url" usage:"specify the base url of the jenkins instance we're indexing"`
	Jobs    *cli.StringSlice `json:"job"      usage:"provide an initial list of jobs to analyze"`
}
```


#### type GetJobResponse

```go
type GetJobResponse struct {
	Builds []Build `json:"builds"`
}
```


#### type GetPipelineResponse

```go
type GetPipelineResponse struct {
	ID                  string          `json:"id"`
	Status              string          `json:"status"`
	StartTimeMillis     int64           `json:"startTimeMillis"`
	DurationMillis      int64           `json:"durationMillis"`
	QueueDurationMillis int64           `json:"queueDurationMillis"`
	PauseDurationMillis int64           `json:"pauseDurationMillis"`
	Stages              []PipelineStage `json:"stages"`
}
```


#### type Job

```go
type Job struct {
	Name string `json:"name"`
}
```


#### type Jobs

```go
type Jobs struct {
}
```


#### func (*Jobs) Get

```go
func (jobs *Jobs) Get(ctx context.Context, project string) (*GetJobResponse, error)
```

#### func (*Jobs) List

```go
func (jobs *Jobs) List(ctx context.Context) (*ListJobsResponse, error)
```

#### type ListJobsResponse

```go
type ListJobsResponse struct {
	Jobs []Job `json:"jobs"`
}
```


#### type PipelineStage

```go
type PipelineStage struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Status              string `json:"status"`
	StartTimeMillis     int64  `json:"startTimeMillis"`
	DurationMillis      int64  `json:"durationMillis"`
	PauseDurationMillis int64  `json:"pauseDurationMillis"`
}
```


#### type Pipelines

```go
type Pipelines struct {
}
```


#### func (*Pipelines) Get

```go
func (pipelines *Pipelines) Get(ctx context.Context, project string, build int) (*GetPipelineResponse, error)
```
