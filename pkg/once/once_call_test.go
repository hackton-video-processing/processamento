package once_test

import (
	"sync"
	"testing"

	"github.com/hackton-video-processing/processamento/pkg/once"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnceCallSimultaneousCalls(t *testing.T) {
	// Arrange
	var (
		wg      sync.WaitGroup
		results = make([]int, 10)

		singletonFunc = func() (int, error) {
			return 42, nil
		}
	)

	// Act
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value, err := once.Call[int](singletonFunc)
			require.NoError(t, err)
			results[i] = value
		}(i)
	}

	wg.Wait()

	// Assert
	for i := 1; i < len(results); i++ {
		require.Equal(t, results[0], results[i])
	}
}

func TestOnceCall(t *testing.T) {
	t.Run("singletonFunc returns error", func(t *testing.T) {
		// Arrange
		once.CallFlush()
		expectedError := assert.AnError

		// Act
		_, err := once.Call[int](func() (int, error) {
			return 0, expectedError
		})

		// Assert
		require.ErrorIs(t, err, expectedError)
	})
}
