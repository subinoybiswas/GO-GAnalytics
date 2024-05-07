package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	data "google.golang.org/genproto/googleapis/analytics/data/v1beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

const (
	propertyID             = "391742241"
	serviceAccountJsonPath = "./subinoy-382420-f00fc278b083.json"
)

type Resp struct {
	DimensionValue string `json:"dimension_value"`
	EventCount     string `json:"event_count"`
	ActiveUsers    string `json:"active_users"`
}

func getDataAnalytics(w http.ResponseWriter, r *http.Request) {
	// Load service account credentials from JSON key file
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	creds, err := oauth.NewServiceAccountFromFile(serviceAccountJsonPath, "https://www.googleapis.com/auth/analytics.readonly")
	if err != nil {
		log.Fatalf("Failed to load service account credentials: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create a gRPC connection with credentials
	conn, err := grpc.Dial(
		"analyticsdata.googleapis.com:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithPerRPCCredentials(creds),
	)
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Create a client for the Analytics Data API
	client := data.NewBetaAnalyticsDataClient(conn)

	// Make a request to the GA4 Data API
	request := &data.RunReportRequest{
		Property: "properties/" + propertyID,
		Dimensions: []*data.Dimension{
			{Name: "day"}, // change according to your needs
		},
		Metrics: []*data.Metric{
			{
				Name: "totalUsers",
			},
			{
				Name: "averageSessionDuration",
			},
			{
				Name: "active7DayUsers",
			},
			// Add more metrics as needed
		},
		DateRanges: []*data.DateRange{
			{
				StartDate: "2023-01-21",
				EndDate:   time.Now().Format("2006-01-02"),
			},
		},
	}

	// Record start time
	startTime := time.Now()

	// Execute the function
	response, err := client.RunReport(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to run report: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Record end time
	endTime := time.Now()

	// Print execution time
	fmt.Println("Execution time:", endTime.Sub(startTime))

	// fmt.Println("Report result:")
	// for _, row := range response.Rows {
	// 	fmt.Printf("%s, Event Count: %v, Active Users: %v\n", row.DimensionValues[0].GetValue(), row.MetricValues[0].GetValue(), row.MetricValues[1].GetValue())
	// }
	respData := make([]Resp, len(response.Rows))
	for i, row := range response.Rows {
		respData[i] = Resp{
			DimensionValue: row.DimensionValues[0].GetValue(),
			EventCount:     row.MetricValues[0].GetValue(),
			ActiveUsers:    row.MetricValues[1].GetValue(),
		}
	}

	jsonData, err := json.Marshal(respData)
	if err != nil {
		log.Fatalf("Json Error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

	// Write the response to the HTTP response writer
}

func main() {
	port := ":8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/", getDataAnalytics)
	fmt.Println("Server Listening on PORT", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
