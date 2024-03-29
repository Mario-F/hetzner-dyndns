package providers

import "testing"

type TestStringValueResult struct {
	Inputs []string
	Output string
}

func TestCaptureIPv6(t *testing.T) {
	t.Run("Test capture from string", func(t *testing.T) {
		testings := []TestStringValueResult{
			{
				Inputs: []string{
					"Your IPv6 address on the public Internet appears to be 2001:9e8:e169:4400:23f5:f330:350b:7e80",
					"2001:9e8:e169:4400:23f5:f330:350b:7e80 is your ip address",
					"2001:9e8:e169:4400:23f5:f330:350b:7e80right noise",
					"left noise2001:9e8:e169:4400:23f5:f330:350b:7e80",
					"all2001:9e8:e169:4400:23f5:f330:350b:7e80noise",
				},
				Output: "2001:9e8:e169:4400:23f5:f330:350b:7e80",
			},
			{
				Inputs: []string{
					"2a01:4f8:1c1e:71c9::1",
					"your ipv6 is 2a01:4f8:1c1e:71c9::1, click for...",
					"##     Your IP Address is 2a01:4f8:1c1e:71c9::1 (34850)     ##",
				},
				Output: "2a01:4f8:1c1e:71c9::1",
			},
			{
				Inputs: []string{
					"2001:db8:ea34::71ff:fe0",
					"noiseleft2001:db8:ea34::71ff:fe0 noise right",
				},
				Output: "2001:db8:ea34::71ff:fe0",
			},
			{
				Inputs: []string{
					"fe08::7:8%1",
					"left noise fe08::7:8%1rightnoise",
				},
				Output: "fe08::7:8",
			},
			{
				Inputs: []string{
					"1::3:4:5:6:7:8",
					"asdasd 1::3:4:5:6:7:8noise",
					"<>1::3:4:5:6:7:8-1",
				},
				Output: "1::3:4:5:6:7:8",
			},
		}

		for _, testV6 := range testings {
			for _, testInput := range testV6.Inputs {
				ip, err := captureIPv6(testInput)
				if err != nil {
					t.Error(err)
					t.Errorf("Cant parse string: %v+", testInput)
				}
				if ip != testV6.Output {
					t.Errorf("%s should be %s", ip, testV6.Output)
				}
			}
		}
	})
}

func TestCaptureIPv4(t *testing.T) {
	t.Run("Test capture from string", func(t *testing.T) {
		testings := []TestStringValueResult{
			{
				Inputs: []string{
					"<html><div>you ip is 89.244.207.0</div></html>",
					"89.244.207.0",
				},
				Output: "89.244.207.0",
			},
			{
				Inputs: []string{
					"4.4.4.4",
					"4.4.4.4right noise",
					"leftnoise4.4.4.4",
					"full.4.4.4.4-noise",
				},
				Output: "4.4.4.4",
			},
		}

		for _, testV4 := range testings {
			for _, testInput := range testV4.Inputs {
				ip, err := captureIPv4(testInput)
				if err != nil {
					t.Error(err)
					t.Errorf("Cant parse string: %v+", testInput)
				}
				if ip != testV4.Output {
					t.Errorf("%s should be %s", ip, testV4.Output)
				}
			}
		}
	})
}
