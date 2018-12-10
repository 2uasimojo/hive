// Package aws collects AWS-specific configuration.
package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types/aws"
)

const (
	defaultVPCCIDR = "10.0.0.0/16"
)

var (
	validAWSRegions = map[string]string{
		"ap-northeast-1": "Tokyo",
		"ap-northeast-2": "Seoul",
		"ap-northeast-3": "Osaka-Local",
		"ap-south-1":     "Mumbai",
		"ap-southeast-1": "Singapore",
		"ap-southeast-2": "Sydney",
		"ca-central-1":   "Central",
		"cn-north-1":     "Beijing",
		"cn-northwest-1": "Ningxia",
		"eu-central-1":   "Frankfurt",
		"eu-west-1":      "Ireland",
		"eu-west-2":      "London",
		"eu-west-3":      "Paris",
		"sa-east-1":      "São Paulo",
		"us-east-1":      "N. Virginia",
		"us-east-2":      "Ohio",
		"us-west-1":      "N. California",
		"us-west-2":      "Oregon",
	}
)

// Platform collects AWS-specific configuration.
func Platform() (*aws.Platform, error) {
	longRegions := make([]string, 0, len(validAWSRegions))
	shortRegions := make([]string, 0, len(validAWSRegions))
	for id, location := range validAWSRegions {
		longRegions = append(longRegions, fmt.Sprintf("%s (%s)", id, location))
		shortRegions = append(shortRegions, id)
	}
	regionTransform := survey.TransformString(func(s string) string {
		return strings.SplitN(s, " ", 2)[0]
	})

	defaultRegion := "us-east-1"
	_, ok := validAWSRegions[defaultRegion]
	if !ok {
		panic(fmt.Sprintf("installer bug: invalid default AWS region %q", defaultRegion))
	}

	ssn := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ssn.Config.Credentials = credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})
	_, err := ssn.Config.Credentials.Get()
	if err == credentials.ErrNoValidProvidersFoundInChain {
		err = getCredentials()
		if err != nil {
			return nil, err
		}
	}

	defaultRegionPointer := ssn.Config.Region
	if defaultRegionPointer != nil {
		_, ok := validAWSRegions[*defaultRegionPointer]
		if ok {
			defaultRegion = *defaultRegionPointer
		} else {
			logrus.Warnf("Unrecognized AWS region %q, defaulting to %s", *defaultRegionPointer, defaultRegion)
		}
	}

	sort.Strings(longRegions)
	sort.Strings(shortRegions)
	region, err := asset.GenerateUserProvidedAsset(
		"AWS Region",
		&survey.Question{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The AWS region to be used for installation.",
				Default: fmt.Sprintf("%s (%s)", defaultRegion, validAWSRegions[defaultRegion]),
				Options: longRegions,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := regionTransform(ans).(string)
				i := sort.SearchStrings(shortRegions, choice)
				if i == len(shortRegions) || shortRegions[i] != choice {
					return errors.Errorf("invalid region %q", choice)
				}
				return nil
			}),
			Transform: regionTransform,
		},
		"OPENSHIFT_INSTALL_AWS_REGION",
	)
	if err != nil {
		return nil, err
	}

	userTags := map[string]string{}
	if value, ok := os.LookupEnv("_CI_ONLY_STAY_AWAY_OPENSHIFT_INSTALL_AWS_USER_TAGS"); ok {
		if err := json.Unmarshal([]byte(value), &userTags); err != nil {
			return nil, errors.Wrapf(err, "_CI_ONLY_STAY_AWAY_OPENSHIFT_INSTALL_AWS_USER_TAGS contains invalid JSON: %s", value)
		}
	}

	return &aws.Platform{
		VPCCIDRBlock: defaultVPCCIDR,
		Region:       region,
		UserTags:     userTags,
	}, nil
}

func getCredentials() error {
	keyID, err := asset.GenerateUserProvidedAsset(
		"AWS Access Key ID",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "AWS Access Key ID",
				Help:    "The AWS access key ID to use for installation (this is not your username).\nhttps://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html",
			},
		},
		"",
	)
	if err != nil {
		return err
	}

	secretKey, err := asset.GenerateUserProvidedAsset(
		"AWS Access Key ID",
		&survey.Question{
			Prompt: &survey.Password{
				Message: "AWS Secret Access Key",
				Help:    "The AWS secret access key corresponding to your access key ID (this is not your password).",
			},
		},
		"",
	)
	if err != nil {
		return err
	}

	tmpl, err := template.New("aws-credentials").Parse(`# Created by openshift-install
# https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
[default]
aws_access_key_id={{.KeyID}}
aws_secret_access_key={{.SecretKey}}
`)
	if err != nil {
		return err
	}

	path := defaults.SharedCredentialsFilename()
	logrus.Infof("Writing AWS credentials to %q (https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)", path)
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, map[string]string{
		"KeyID":     keyID,
		"SecretKey": secretKey,
	})
}
