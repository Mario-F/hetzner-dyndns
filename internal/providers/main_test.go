package providers

import "testing"

func TestCaptureIPv6(t *testing.T) {
	t.Run("Test capture from string", func(t *testing.T) {
		testString := "Your IPv6 address on the public Internet appears to be 2001:9e8:e169:4400:23f5:f330:350b:7e80"
		testStringResult := "2001:9e8:e169:4400:23f5:f330:350b:7e80"
		ip, err := captureIPv6(testString)
		if err != nil {
			t.Error(err)
			t.Errorf("Cant parse string: %v+", testString)
		}
		if ip != testStringResult {
			t.Errorf("%s should be %s", ip, testStringResult)
		}
	})
}

func TestCaptureIPv4(t *testing.T) {
	t.Run("Test capture from string", func(t *testing.T) {
		testString := "<html><div>you ip is 89.244.207.0</div></html>"
		testStringResult := "89.244.207.0"
		ip, err := captureIPv4(testString)
		if err != nil {
			t.Error(err)
			t.Errorf("Cant parse string: %v+", testString)
		}
		if ip != testStringResult {
			t.Errorf("%s should be %s", ip, testStringResult)
		}
	})
}
