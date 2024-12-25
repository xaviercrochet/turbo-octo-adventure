package app

import (
	"sync"
	"testing"
)

func TestSetAndgetSelectedUsername(t *testing.T) {
	initial := getSelectedUsername()
	defer setSelectedUsername(initial)

	tests := []struct {
		name     string
		username string
		expected string
	}{
		{
			name:     "basic set and get",
			username: "Test Username",
			expected: "Test Username",
		},
		{
			name:     "empty username",
			username: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setSelectedUsername(tt.username)
			result := getSelectedUsername()
			if result != tt.expected {
				t.Errorf("getSelectedUsername() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

/*


  Verify that no race conditions are triggered if run concurrently
  This test has to be run with -race flag or test-race target
*/

func TestRaceCondition(t *testing.T) {
	initial := getSelectedUsername()
	defer setSelectedUsername(initial)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			setSelectedUsername("user1")
		}()

		go func() {
			defer wg.Done()
			_ = getSelectedUsername()
		}()
	}

	wg.Wait()
}
