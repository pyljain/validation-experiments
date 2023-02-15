package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const (
	nodejsOTELPackage = "@opentelemetry"
	goOTELPackage     = "go.opentelemetry.io/otel"
	apiRoute          = "/api/"
	restOpenApiSpec3  = "openapi: \"3"
)

func main() {
	// Accept user input for root directory path
	rootDir := os.Args[1]
	if rootDir == "" {
		fmt.Println("Please provide the root directory path to scan")
	}

	// Create resultset
	findings := []Finding{}
	// Traverse directory tree structure
	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "node_modules" {
			return filepath.SkipDir
		}

		// Look for presence of otel libraries in Node
		if d.Name() == "package.json" {
			result, err := FindStringInFile(nodejsOTELPackage, path)
			if err != nil {
				return err
			}

			findings = append(findings, Finding{
				Category:    FindingCategoryTelemetryStandard,
				Description: "Looking for the presence of the OTEL library",
				Location:    path,
				Compliant:   result,
			})
		}

		// Look for presence of otel libraries in Node
		if d.Name() == "go.mod" {
			result, err := FindStringInFile(goOTELPackage, path)
			if err != nil {
				return err
			}

			findings = append(findings, Finding{
				Category:    FindingCategoryTelemetryStandard,
				Description: "Looking for the presence of the OTEL library",
				Location:    path,
				Compliant:   result,
			})

		}

		if d.Name() == "go.mod" || d.Name() == "package.json" {
			parentDirectory := filepath.Dir(path)
			finding, err := EvaluatePresenceOfAPIStandard(parentDirectory)
			if err != nil {
				return err
			}

			if finding != nil {
				findings = append(findings, *finding)
			}
		}

		return nil
	})

	// Print findings
	PrintResults(findings)
}

func FindStringInFile(s string, filepath string) (bool, error) {
	contentBytes, err := os.ReadFile(filepath)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(contentBytes), s), nil
}

func EvaluatePresenceOfAPIStandard(serviceRoot string) (*Finding, error) {
	foundAPIRoute := false
	foundAPISpec := false

	filepath.WalkDir(serviceRoot, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "node_modules" {
			return filepath.SkipDir
		}

		// Look for presence of otel libraries in Node
		if filepath.Ext(path) == ".go" || filepath.Ext(path) == ".js" {
			var err error
			foundAPIRoute, err = FindStringInFile(apiRoute, path)
			if err != nil {
				return err
			}
		}

		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			var err error
			foundAPISpec, err = FindStringInFile(restOpenApiSpec3, path)
			if err != nil {
				return err
			}
		}

		return nil

	})

	f := &Finding{
		Category:    FindingCategoryRESTAPIStandard,
		Description: "Evaluating presence of an OpenAPI spec for REST APIs",
		Location:    serviceRoot,
		Compliant:   (foundAPISpec && foundAPIRoute) || foundAPISpec || (!foundAPISpec && !foundAPIRoute),
	}
	return f, nil
}

func PrintResults(findings []Finding) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Category", "Description", "Location", "Outcome")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, finding := range findings {

		outcome := color.New(color.FgGreen).Sprint("\xE2\x9C\x94")
		if !finding.Compliant {
			outcome = color.New(color.FgRed).Sprint("\u2A09")
		}

		tbl.AddRow(finding.Category, finding.Description, finding.Location, outcome)
	}

	tbl.Print()
}
