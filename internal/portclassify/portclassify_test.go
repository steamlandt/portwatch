package portclassify_test

import (
	"testing"

	"github.com/user/portwatch/internal/portclassify"
)

func makePort(number int, proto string) portclassify.Port {
	return portclassify.Port{Number: number, Protocol: proto}
}

func TestClassifySystemPort(t *testing.T) {
	c := portclassify.New()
	r := c.Classify(makePort(80, "tcp"))
	if r.Class != portclassify.ClassSystem {
		t.Fatalf("expected system, got %s", r.Class)
	}
}

func TestClassifyRegisteredPort(t *testing.T) {
	c := portclassify.New()
	r := c.Classify(makePort(8080, "tcp"))
	if r.Class != portclassify.ClassRegistered {
		t.Fatalf("expected registered, got %s", r.Class)
	}
}

func TestClassifyDynamicPort(t *testing.T) {
	c := portclassify.New()
	r := c.Classify(makePort(55000, "udp"))
	if r.Class != portclassify.ClassDynamic {
		t.Fatalf("expected dynamic, got %s", r.Class)
	}
}

func TestClassifyBoundaryPorts(t *testing.T) {
	c := portclassify.New()
	cases := []struct {
		number int
		want   portclassify.Class
	}{
		{0, portclassify.ClassSystem},
		{1023, portclassify.ClassSystem},
		{1024, portclassify.ClassRegistered},
		{49151, portclassify.ClassRegistered},
		{49152, portclassify.ClassDynamic},
		{65535, portclassify.ClassDynamic},
	}
	for _, tc := range cases {
		r := c.Classify(makePort(tc.number, "tcp"))
		if r.Class != tc.want {
			t.Errorf("port %d: expected %s, got %s", tc.number, tc.want, r.Class)
		}
	}
}

func TestClassifyAllReturnsMapKeyedByNumber(t *testing.T) {
	c := portclassify.New()
	ports := []portclassify.Port{
		makePort(22, "tcp"),
		makePort(3306, "tcp"),
		makePort(60000, "udp"),
	}
	m := c.ClassifyAll(ports)
	if len(m) != 3 {
		t.Fatalf("expected 3 results, got %d", len(m))
	}
	if m[22].Class != portclassify.ClassSystem {
		t.Errorf("port 22 should be system")
	}
	if m[3306].Class != portclassify.ClassRegistered {
		t.Errorf("port 3306 should be registered")
	}
	if m[60000].Class != portclassify.ClassDynamic {
		t.Errorf("port 60000 should be dynamic")
	}
}

func TestResultLabelIncludesProtocol(t *testing.T) {
	c := portclassify.New()
	r := c.Classify(makePort(443, "tcp"))
	if r.Label != "system/tcp" {
		t.Errorf("unexpected label: %s", r.Label)
	}
}
