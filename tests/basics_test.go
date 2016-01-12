package cacheorganizer_test

import "time"
import "testing"
import "github.com/stretchr/testify/assert"

func TestBasic(t *testing.T) {
    testValues := []int{1, 2, 3, 4}
    result := runBasicTest(testValues, 0, []time.Duration{time.Minute, time.Minute, time.Minute, time.Minute}, []time.Duration{0, 0, 0, 0})
    
    assert.NotNil(t, result, "Result is valid.")
    assert.Nil(t, result.Error, "Result does not contain error.")
    assert.Equal(t, 1, result.Value.(int), "Result is 1.")
    assert.Equal(t, 1, testValues[0], "Cache value is 1.")
}

func TestOverwriteCache(t *testing.T) {
    testValues := []int{1, 2, 3, 4}
    result := runBasicTest(testValues, 1, []time.Duration{time.Minute, time.Minute, time.Minute, time.Minute}, []time.Duration{0, 0, 0, 0})
    
    assert.NotNil(t, result, "Result is valid.")
    assert.Nil(t, result.Error, "Result does not contain error.")
    assert.Equal(t, 2, result.Value.(int), "Result is 2.")
    assert.Equal(t, 2, testValues[0], "Cache value is 2.")
    assert.Equal(t, 2, testValues[1], "Second-level cache value is 2.")
}

func TestTimeoutOverwriteCache(t *testing.T) {
    testValues := []int{1, 2, 3, 4}
    result := runBasicTest(testValues, 0, []time.Duration{50 * time.Millisecond, time.Minute, time.Minute, time.Minute}, []time.Duration{100 * time.Millisecond, 0, 0, 0})
    
    assert.NotNil(t, result, "Result is valid.")
    assert.Nil(t, result.Error, "Result does not contain error.")
    assert.Equal(t, 2, result.Value.(int), "Result is 2.")
    assert.Equal(t, 2, testValues[0], "Cache value is 2.")
    assert.Equal(t, 2, testValues[1], "Second-level cache value is 2.")
}

func TestRespondAfterTimeout(t *testing.T) {
    testValues := []int{1, 2, 3, 4}
    result := runBasicTest(testValues, 0, []time.Duration{50 * time.Millisecond, time.Minute, time.Minute, time.Minute}, []time.Duration{100 * time.Millisecond, time.Minute, 0, 0})
    
    assert.NotNil(t, result, "Result is valid.")
    assert.Nil(t, result.Error, "Result does nto contain error.")
    assert.Equal(t, 1, result.Value.(int), "Result is 1.")
    assert.Equal(t, 1, testValues[0], "Cache value is 1.")
}

func TestThirdLevelResponse(t *testing.T) {
    testValues := []int{1, 2, 3, 4}
    result := runBasicTest(testValues, 1, []time.Duration{time.Minute, 50 * time.Millisecond, time.Minute, time.Minute}, []time.Duration{0, 100 * time.Millisecond, 0, 0})
    
    assert.NotNil(t, result, "Result is valid.")
    assert.Nil(t, result.Error, "Result does not contain error.")
    assert.Equal(t, 3, result.Value.(int), "Result is 3.")
    assert.Equal(t, 3, testValues[0], "Cache value is 3.")
    assert.Equal(t, 3, testValues[1], "Second-level cache value is 3.")
    assert.Equal(t, 3, testValues[2], "Third-level cache value is 3.")
    assert.Equal(t, 4, testValues[3], "Fourth-level cache is 4.")
}