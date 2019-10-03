/*
   Define and read the configuration file containg the desired branch
   protection policy.
*/
package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// BranchPolicy represents the config file detailing the desired policy.
type BranchPolicy struct {
	*RequireStatusChecks        `yaml:"requireStatusChecks"`
	*RequiredPullRequestReviews `yaml:"requiredPullRequestReviews"`
	RequireSignatures           bool     `yaml:"requireSignatures"`
	EnforceAdmins               bool     `yaml:"enforceAdmins"`
	NotifyUsers                 []string `yaml:"notifyUsers,omitempty"`
}

type RequireStatusChecks struct {
	Strict   bool     `yaml:"strict"`
	Contexts []string `yaml:"contexts"`
}

type RequiredPullRequestReviews struct {
	DismissStaleReviews          bool `yaml:"dismissStaleReviews"`
	RequiredApprovingReviewCount int  `yaml:"requiredApprovingReviewCount"`
}

func (b *BranchPolicy) String() string {
	return fmt.Sprintf(
		"BranchPolicy{RequireStatusChecks:%+v, RequiredPullRequestReviews:%+v, RequireSignatures:%v, EnforceAdmins:%v, NotifyUsers:%v}",
		b.RequireStatusChecks, b.RequiredPullRequestReviews, b.RequireSignatures, b.EnforceAdmins, b.NotifyUsers,
	)
}

// UnmarshalYAML allows unmarshalling and application of default values.
func (b *BranchPolicy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type branchPolicyType BranchPolicy
	defaultBranchPolicy := branchPolicyType{
		RequireStatusChecks: &RequireStatusChecks{
			Strict:   true,
			Contexts: []string{},
		},
		RequiredPullRequestReviews: &RequiredPullRequestReviews{
			DismissStaleReviews:          true,
			RequiredApprovingReviewCount: 2,
		},
		RequireSignatures: true,
		EnforceAdmins:     true,
		NotifyUsers:       []string{},
	}
	if err := unmarshal(&defaultBranchPolicy); err != nil {
		return err
	}

	*b = BranchPolicy(defaultBranchPolicy)
	return nil
}

// readBranchProtectionConfig returns a BranchPolicy object with the
// information from the configuration file.
func readBranchProtectionConfig(config string) (*BranchPolicy, error) {
	yamlFile, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	defaultBranchPolicy := &BranchPolicy{}
	if err := yaml.Unmarshal(yamlFile, defaultBranchPolicy); err != nil {
		return nil, err
	}

	return defaultBranchPolicy, nil
}
