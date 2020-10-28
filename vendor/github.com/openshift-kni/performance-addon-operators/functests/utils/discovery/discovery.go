package discovery

import (
	"context"
	"fmt"
	"os"
	"strconv"

	performancev2 "github.com/openshift-kni/performance-addon-operators/api/v2"
	testclient "github.com/openshift-kni/performance-addon-operators/functests/utils/client"
	"github.com/openshift-kni/performance-addon-operators/functests/utils/profiles"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ErrProfileNotFound = fmt.Errorf("Profile not found in discovery mode")

// ConditionIterator is the function that accepts element of a PerformanceProfile and returns boolean
type ConditionIterator func(performancev2.PerformanceProfile) bool

// Enabled indicates whether test discovery mode is enabled.
func Enabled() bool {
	discoveryMode, _ := strconv.ParseBool(os.Getenv("DISCOVERY_MODE"))
	return discoveryMode
}

// GetDiscoveryPerformanceProfile returns an existing profile in the cluster with the most nodes using it.
// In case no profile exists - return nil
func GetDiscoveryPerformanceProfile() (*performancev2.PerformanceProfile, error) {
	performanceProfiles, err := profiles.All()
	if err != nil {
		return nil, err
	}
	return getDiscoveryPerformanceProfile(performanceProfiles.Items)
}

// GetFilteredDiscoveryPerformanceProfile returns an existing profile in the cluster with the most nodes using it
// from a a filtered profiles list by the filter function passed as an argument.
// In case no profile exists - return nil
func GetFilteredDiscoveryPerformanceProfile(iterator ConditionIterator) (*performancev2.PerformanceProfile, error) {
	performanceProfiles, err := profiles.All()
	if err != nil {
		return nil, err
	}
	return getDiscoveryPerformanceProfile(filter(performanceProfiles.Items, iterator))
}

func getDiscoveryPerformanceProfile(performanceProfiles []performancev2.PerformanceProfile) (*performancev2.PerformanceProfile, error) {
	var currentProfile *performancev2.PerformanceProfile = nil
	maxNodesNumber := 0
	for _, profile := range performanceProfiles {
		selector := labels.SelectorFromSet(profile.Spec.NodeSelector)

		profileNodes := &corev1.NodeList{}
		if err := testclient.Client.List(context.TODO(), profileNodes, &client.ListOptions{LabelSelector: selector}); err != nil {
			return nil, err
		}

		if len(profileNodes.Items) > maxNodesNumber {
			currentProfile = &profile
			maxNodesNumber = len(profileNodes.Items)
		}
	}

	if currentProfile == nil {
		return nil, ErrProfileNotFound
	}
	return currentProfile, nil
}

func filter(performanceProfiles []performancev2.PerformanceProfile, iterator ConditionIterator) []performancev2.PerformanceProfile {
	var result = make([]performancev2.PerformanceProfile, 0)
	for _, profile := range performanceProfiles {
		if iterator(profile) {
			result = append(result, profile)
		}
	}
	return result
}
