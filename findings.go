package main

type Finding struct {
	Category    FindingCategory
	Description string
	Location    string
	Compliant   bool
}

type FindingCategory string

const (
	FindingCategoryTelemetryStandard FindingCategory = "Telemetry standard"
	FindingCategoryRESTAPIStandard   FindingCategory = "REST API standard"
)
