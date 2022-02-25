package providers

import "testing"

type TestStringValueResult struct {
	Input  string
	Output string
}

func TestCaptureIPv6(t *testing.T) {
	t.Run("Test capture from string", func(t *testing.T) {
		testings := []TestStringValueResult{
			{
				Input:  "Your IPv6 address on the public Internet appears to be 2001:9e8:e169:4400:23f5:f330:350b:7e80",
				Output: "2001:9e8:e169:4400:23f5:f330:350b:7e80",
			},
			{
				Input:  "2001:9e8:e169:4400:23f5:f330:350b:7e80 is your ip address",
				Output: "2001:9e8:e169:4400:23f5:f330:350b:7e80",
			},
		}

		for _, testV6 := range testings {
			ip, err := captureIPv6(testV6.Input)
			if err != nil {
				t.Error(err)
				t.Errorf("Cant parse string: %v+", testV6.Input)
			}
			if ip != testV6.Output {
				t.Errorf("%s should be %s", ip, testV6.Output)
			}
		}
	})
}

func TestCaptureIPv4(t *testing.T) {
	t.Run("Test capture from string", func(t *testing.T) {
		testings := []TestStringValueResult{
			{
				Input:  "<html><div>you ip is 89.244.207.0</div></html>",
				Output: "89.244.207.0",
			},
			{
				Input:  "89.244.207.0",
				Output: "89.244.207.0",
			},
		}

		for _, testV4 := range testings {
			ip, err := captureIPv4(testV4.Input)
			if err != nil {
				t.Error(err)
				t.Errorf("Cant parse string: %v+", testV4.Input)
			}
			if ip != testV4.Output {
				t.Errorf("%s should be %s", ip, testV4.Output)
			}
		}
	})
}
