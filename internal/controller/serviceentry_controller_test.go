package controller

import (
	"testing"
)

func TestCreateDefaultWeightForNewEndpoints(t *testing.T) {
	tests := []struct {
		name                         string
		existingWeights              map[string]uint32
		NewEndpointsPercentileWeight int
		MinimumWeight                int
		MaximumWeight                int
		expectedWeight               uint32
	}{
		{
			name:                         "No existing weights",
			existingWeights:              map[string]uint32{},
			NewEndpointsPercentileWeight: 50, // Doesn't matter when existingWeights is empty
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               50, // (0 + 100) / 2
		},
		{
			name:                         "Existing weights with 50% percentile",
			existingWeights:              map[string]uint32{"ip1": 10, "ip2": 20, "ip3": 30},
			NewEndpointsPercentileWeight: 50,
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               10, // Average of lowest 1 weight (n=1)
		},
		{
			name:                         "Existing weights with 80% percentile",
			existingWeights:              map[string]uint32{"ip1": 10, "ip2": 20, "ip3": 30, "ip4": 40, "ip5": 50},
			NewEndpointsPercentileWeight: 80,
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               25, // Average of lowest 4 weights (n=4)
		},
		{
			name:                         "Existing weights with 0% percentile",
			existingWeights:              map[string]uint32{"ip1": 15, "ip2": 25, "ip3": 35},
			NewEndpointsPercentileWeight: 0,
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               15, // Average of lowest 1 weight (n=1)
		},
		{
			name:                         "Existing weights with 100% percentile",
			existingWeights:              map[string]uint32{"ip1": 10, "ip2": 20, "ip3": 30},
			NewEndpointsPercentileWeight: 100,
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               20, // Average of all weights (n=3)
		},
		{
			name:                         "Minimum and Maximum weights adjusted",
			existingWeights:              map[string]uint32{},
			NewEndpointsPercentileWeight: 50, // Doesn't matter when existingWeights is empty
			MinimumWeight:                10,
			MaximumWeight:                90,
			expectedWeight:               50, // (10 + 90) / 2
		},
		{
			name:                         "Existing weights with 33% percentile",
			existingWeights:              map[string]uint32{"ip1": 5, "ip2": 15, "ip3": 25, "ip4": 35, "ip5": 45},
			NewEndpointsPercentileWeight: 33,
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               5, // Average of lowest 1 weight (n=1)
		},
		{
			name:                         "Existing weights with less than 1 endpoint after percentile calculation",
			existingWeights:              map[string]uint32{"ip1": 50},
			NewEndpointsPercentileWeight: 25,
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               50, // Should consider at least one endpoint
		},
		{
			name:                         "Existing weights with high percentile but few endpoints",
			existingWeights:              map[string]uint32{"ip1": 10, "ip2": 20},
			NewEndpointsPercentileWeight: 90,
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               10,
		},
		{
			name:                         "Existing weights with 10% percentile",
			existingWeights:              map[string]uint32{"ip1": 5, "ip2": 15, "ip3": 25, "ip4": 35, "ip5": 45},
			NewEndpointsPercentileWeight: 1,
			MinimumWeight:                0,
			MaximumWeight:                100,
			expectedWeight:               5, // Average of lowest 1 weight (n=1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ServiceEntryReconciler{
				NewEndpointsPercentileWeight: tt.NewEndpointsPercentileWeight,
				MinimumWeight:                tt.MinimumWeight,
				MaximumWeight:                tt.MaximumWeight,
			}
			got := r.CreateDefaultWeightForNewEndpoints(tt.existingWeights)
			if got != tt.expectedWeight {
				t.Errorf("CreateDefaultWeightForNewEndpoints() = %v, want %v", got, tt.expectedWeight)
			}
		})
	}
}
