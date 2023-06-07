package gcputil

import (
	"bytes"
	"context"
	"github.com/secureworks/errors"
	"golang.org/x/oauth2/google"
	"os"
	"os/exec"
	"strings"
)

func GetGoogleCloudCLIProject(ctx context.Context) (string, error) {
	out := bytes.Buffer{}
	getProjectCmd := exec.CommandContext(ctx, "gcloud", "config", "get", "project")
	getProjectCmd.Stderr = os.Stderr
	getProjectCmd.Stdout = &out
	if err := getProjectCmd.Run(); err != nil {
		return "", errors.Chain(err, "failed to get GCP project from gcloud CLI")
	} else {
		return strings.TrimSpace(out.String()), nil
	}
}

func GetADCProject(ctx context.Context) (string, error) {
	credentials, adcErr := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/compute")
	if adcErr != nil {
		return "", errors.Chain(adcErr, "failed to get GCP project from ADC")
	}

	gcpProjectID := credentials.ProjectID
	if gcpProjectID != "" {
		return gcpProjectID, nil
	}

	return "", errors.NewWithStackTrace("ADC project is empty")
}

func GetDefaultProjectID(ctx context.Context) (string, error) {
	project, adcErr := GetADCProject(ctx)
	if adcErr != nil {
		project, cliErr := GetGoogleCloudCLIProject(ctx)
		if cliErr != nil {
			err := errors.NewMultiError(adcErr, cliErr)
			return "", errors.Chain(err, "failed to get GCP project")
		} else {
			return project, nil
		}
	} else {
		return project, nil
	}
}
