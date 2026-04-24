package main

import (
	"bytes"
	"compress/bzip2"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/deutscherwetterdienst/dwd-downloader-go/internal/logger"
	"github.com/deutscherwetterdienst/dwd-downloader-go/internal/models"
	"github.com/deutscherwetterdienst/dwd-downloader-go/internal/version"
	"github.com/urfave/cli/v2"
)

// loadModels loads model configurations from models.json
func loadModels() models.Available {
	// Load models.json
	data, err := os.ReadFile(filepath.Join("internal", "models", "models.json"))
	if err != nil {
		log.Fatalf("Failed to read models.json: %v", err)
	}

	var modelsList []models.ModelConfig
	if err := json.Unmarshal(data, &modelsList); err != nil {
		log.Fatalf("Failed to parse models.json: %v", err)
	}

	available := models.Available{
		Models: make(map[string]models.ModelConfig),
		Grids:  make(map[string]string),
	}

	for _, model := range modelsList {
		available.Models[model.Model] = model
		for _, grid := range model.Grids {
			available.Grids[grid] = grid
		}
	}

	return available
}

func getGribFileURL(model, grid, param string, timestep int, timestamp time.Time, models models.Available) string {
	cfg, ok := models.Models[model]
	if !ok {
		logger.Logger.Fatalf("Unknown model: %s", model)
	}

	if grid == "" {
		logger.Logger.Printf("No grid specified. Trying to use default.")
		if len(cfg.Grids) > 0 {
			grid = cfg.Grids[0]
			logger.Logger.Printf("Grid type '%s' selected", grid)
		} else {
			logger.Logger.Fatalf("No grids available for model %s", model)
		}
	} else if _, ok := models.Grids[grid]; !ok {
		logger.Logger.Printf("Unknown grid type '%s' for model '%s'.", grid, model)
	}

	levtype := "single-level"
	modelrun := timestamp.Hour()
	timestampStr := timestamp.Format("20060102")

	url := cfg.Pattern.SingleLevel
	// Replace placeholders in the pattern
	url = strings.ReplaceAll(url, "{model!L}", strings.ToLower(model))
	url = strings.ReplaceAll(url, "{model}", model)
	url = strings.ReplaceAll(url, "{param!L}", strings.ToLower(param))
	url = strings.ReplaceAll(url, "{param!U}", strings.ToUpper(param))
	url = strings.ReplaceAll(url, "{grid}", grid)
	url = strings.ReplaceAll(url, "{scope}", cfg.Scope)
	url = strings.ReplaceAll(url, "{levtype}", levtype)
	url = strings.ReplaceAll(url, "{modelrun:>02d}", fmt.Sprintf("%02d", modelrun))
	url = strings.ReplaceAll(url, "{timestamp:%Y%m%d}", timestampStr)
	url = strings.ReplaceAll(url, "{step:>03d}", fmt.Sprintf("%03d", timestep))

	return url
}

func downloadAndExtractBz2FileFromURL(url, destFilePath, destFileName string) error {
	logger.Logger.Printf("downloading file: '%s'", url)

	if destFileName == "" {
		parts := strings.Split(url, "/")
		destFileName = parts[len(parts)-1]
		destFileName = strings.TrimSuffix(destFileName, ".bz2")
	}

	if destFilePath == "" {
		destFilePath = "."
	}

	// Validate destination path to prevent path traversal
	absDest, err := filepath.Abs(destFilePath)
	if err != nil {
		return fmt.Errorf("invalid destination path: %v", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: HTTP status %s", resp.Status)
	}

	compressedData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	binaryData := bzip2.NewReader(bytes.NewReader(compressedData))
	fullFilePath := filepath.Join(absDest, destFileName)
	logger.Logger.Printf("saving file as: '%s'", fullFilePath)

	outFile, err := os.Create(fullFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, binaryData); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	logger.Logger.Println("Done.")
	return nil
}

func downloadGribData(model, grid, param string, minTimeStep, maxTimeStep, timeStepInterval int, timestamp time.Time, destFilePath string, models models.Available, parallel int) error {
	fields := strings.Split(param, ",")

	// Ensure at least 1 parallel
	if parallel < 1 {
		parallel = 1
	}

	// Create a semaphore to limit concurrency
	sem := make(chan struct{}, parallel)

	// Collect errors from all downloads
	var mu sync.Mutex
	var errors []string

	// Generate all download tasks
	var wg sync.WaitGroup
	for timestep := minTimeStep; timestep <= maxTimeStep; timestep += timeStepInterval {
		for _, field := range fields {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}

			wg.Add(1)
			go func(ts int, f string) {
				defer wg.Done()

				// Acquire semaphore
				sem <- struct{}{}
				defer func() { <-sem }()

				url := getGribFileURL(model, grid, f, ts, timestamp, models)
				if err := downloadAndExtractBz2FileFromURL(url, destFilePath, ""); err != nil {
					mu.Lock()
					errors = append(errors, fmt.Sprintf("%s: %v", url, err))
					mu.Unlock()
				}
			}(timestep, field)
		}
	}

	// Wait for all downloads to complete
	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("some downloads failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func formatDateIso8601(date time.Time) string {
	return date.UTC().Format(time.RFC3339)
}

func getTimestampString(date time.Time) string {
	modelrun := fmt.Sprintf("%02d", date.Hour())
	return date.Format("20060102") + modelrun
}

func main() {
	app := &cli.App{
		Name:    "downloader",
		Usage:   "Downloads NWP model data in GRIB2 format from DWD's Open Data file server",
		Version: version.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "model",
				Usage:    "the NWP model name",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "grid",
				Usage:    "the model grid",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "single-level-fields",
				Usage:    "one or more single-level model fields that should be downloaded, e.g. t_2m,tmax_2m,clch,pmsl, ...",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "min-time-step",
				Usage: "the minimum forecast time step to download (default=0)",
				Value: 0,
			},
			&cli.IntFlag{
				Name:  "max-time-step",
				Usage: "the maximum forecast time step to download, e.g. 12 will download time steps from min-time-step - 12 (default=0)",
				Value: 0,
			},
			&cli.IntFlag{
				Name:  "time-step-interval",
				Usage: "the interval (in hours) between forecast time steps to download. (Default=1)",
				Value: 1,
			},
			&cli.StringFlag{
				Name:     "timestamp",
				Usage:    "the time stamp of the dataset, e.g. '2020-06-26 18:00'. Uses latest available if no timestamp is specified.",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "directory",
				Usage:    "the download directory, defaults to working directory",
				Required: false,
			},
			&cli.IntFlag{
				Name:  "parallel",
				Usage: "number of parallel downloads (default=1, use >1 to enable concurrent downloads)",
				Value: 1,
			},
		},
		Action: func(c *cli.Context) error {
			model := c.String("model")
			grid := c.String("grid")
			fields := c.String("single-level-fields")
			minTimeStep := c.Int("min-time-step")
			maxTimeStep := c.Int("max-time-step")
			timeStepInterval := c.Int("time-step-interval")
			timestampStr := c.String("timestamp")
			directory := c.String("directory")
			parallel := c.Int("parallel")

			if directory == "" {
				directory = "."
			}

			// Load model configurations
			modelsAvail := loadModels()

			var timestamp time.Time
			if timestampStr != "" {
				// Try multiple timestamp formats
				formats := []string{
					time.RFC3339,
					"2006-01-02 15:04:05",
					"2006-01-02 15:04",
					"2006-01-02",
				}
				var parseErr error
				for _, format := range formats {
					timestamp, parseErr = time.Parse(format, timestampStr)
					if parseErr == nil {
						break
					}
				}
				if parseErr != nil {
					return fmt.Errorf("failed to parse timestamp '%s': %v (supported formats: RFC3339, YYYY-MM-DD HH:MM:SS, YYYY-MM-DD HH:MM, YYYY-MM-DD)", timestampStr, parseErr)
				}
			} else {
				if model == "" {
					// Pick any available model
					for m := range modelsAvail.Models {
						model = m
						break
					}
				}
				if model == "" {
					return fmt.Errorf("no model specified and no models available")
				}

				modelCfg, ok := modelsAvail.Models[model]
				if !ok {
					return fmt.Errorf("unknown model: %s", model)
				}
				timestamp = models.GetMostRecentModelTimestamp(modelCfg)
			}

			if grid == "" {
				if model != "" {
					modelCfg, ok := modelsAvail.Models[model]
					if ok && len(modelCfg.Grids) > 0 {
						grid = modelCfg.Grids[0]
					}
				}
			}

			logger.Logger.Printf(`
---------------
Model: %s
Grid: %s
Fields: %s
Minimum time step: %d
Maximum time step: %d
Time step interval: %d
Timestamp: %s
Model run: %02d
Destination: %s
---------------
`,
				model,
				grid,
				fields,
				minTimeStep,
				maxTimeStep,
				timeStepInterval,
				timestamp.Format("2006-01-02"),
				timestamp.Hour(),
				directory,
			)

			if err := downloadGribData(model, grid, fields, minTimeStep, maxTimeStep, timeStepInterval, timestamp, directory, modelsAvail, parallel); err != nil {
				return fmt.Errorf("download failed: %v", err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
