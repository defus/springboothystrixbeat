package beater

import(
    "fmt"
    "time"
    "strings"
    "regexp"
    "crypto/tls"

    "net/http"
    "io/ioutil"
    "encoding/json"

    "github.com/elastic/beats/libbeat/beat"
    "github.com/elastic/beats/libbeat/logp"
    "github.com/elastic/beats/libbeat/common"
)

//
// Http Response Struct for Concurrent GET's
//
type HttpResponse struct {
    Status string
    Body   []byte
    Url    string
}

//
// Metrics List
//
type MetricsList struct {
    Names []string `json: "names"`
}

//
// Metric Endpoint Response
//
type Metric struct {
    Name string                `json: "statistic"`
    Measurements []Measurement `json: "measurements"`
}

//
// Sub Struct conatining micrometer.io Measurement struct
//
type Measurement struct {
    Statistic string `json: "statistic"`
    Value float32    `json: "value"`
}

//
// Perform a HTTP GET and put the response back on the channel
//
// @param url
// @param ch
//
func DoHttpGet(url string, insecureSkipVerify bool) (*HttpResponse, error) {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
    }

    client := http.Client{
        Timeout: time.Duration(10 * time.Second),
	Transport: tr,
    }
    request, _ := http.NewRequest("GET", url, nil)
    // TODO -- add basic auth
//    request.SetBasicAuth(bt.config.Username, bt.config.Password)
    response, err := client.Do(request)
    if err != nil {
	    return nil, err
    }
    defer response.Body.Close()

    if response.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP%s", response.Status)
    }

    // Fetch "good" Response Body
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
	return nil, err
    }

    return &HttpResponse{response.Status, body, url}, nil
}

//
// Process the Metrics Actuator
//
// @access /actuator/metrics
//
func (bt *Springboothystrixbeat) ProcessMetricsActuator(b *beat.Beat) {
    url := fmt.Sprintf("%s/actuator/metrics", bt.config.Host)

    // Fetch all available Metric Endpoints
    response, err := DoHttpGet(url, bt.config.InsecureSkipVerify)
    if err != nil {
	    logp.Err("Error loading metrics: %v", err)
	    return
    }

    list     := MetricsList{}
    json.Unmarshal(response.Body, &list)
    logp.Debug("springboothystrixbeat", "Liste des métriques : %v", fmt.Sprintf("%v", list))

    // Commit Metric Requests
    requests := 0
    metricResponses := make([]*HttpResponse, len(list.Names))

    for i, metric := range(list.Names) {
	var response HttpResponse
        if existsIn(metric, bt.config.Include) == true {
	    res, err := DoHttpGet(fmt.Sprintf("%s/%s", url, metric), bt.config.InsecureSkipVerify)
	    if err == nil {
		response = *res
                requests++
	    }else{
		logp.Err("Error loading metric value of %v", metric)
	    }
        }
	metricResponses[i] = &response
    }

    // Init Metrics Map
    fields := common.MapStr{
        "type": b.Info.Name,
    }
    logp.Debug("springboothystrixbeat", fmt.Sprintf("Nombre total de métriques actuator disponbles: %d - Nombre de valeur de métriques actuator recupérées: %d", len(list.Names), requests))

    // Join Results to Event
    expr := regexp.MustCompile(`(\.|:)`)
    for it := 0; it < len(list.Names); it++ {
        response := metricResponses[it]
	if response != nil {
            metric   := Metric{}
            json.Unmarshal(response.Body, &metric)
            prefix := expr.ReplaceAllString(metric.Name, "_")
            for _, measurement := range(metric.Measurements) {
                key := fmt.Sprintf("%s_%s", prefix, strings.ToLower(measurement.Statistic))
                fields[key] = measurement.Value
            }
	}
    }

    // Push Event
    event := beat.Event{
        Timestamp: time.Now(),
        Fields: fields,
    }
    bt.client.Publish(event)
    logp.Info(fmt.Sprintf("Event with %d metrics sent.", len(fields)))
}
