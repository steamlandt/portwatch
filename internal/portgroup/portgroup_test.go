package portgroup_test

import (
	"testing"

	"github.com/user/portwatch/internal/portgroup"
)

func TestCategorizeKnownWebPort(t *testing.T) {
	g := portgroup.New(nil)
	if got := g.Categorize(80); got != portgroup.CategoryWeb {
		t.Fatalf("expected web, got %s", got)
	}
}

func TestCategorizeKnownDatabasePort(t *testing.T) {
	g := portgroup.New(nil)
	if got := g.Categorize(5432); got != portgroup.CategoryDatabase {
		t.Fatalf("expected database, got %s", got)
	}
}

func TestCategorizeUnknownPortReturnsOther(t *testing.T) {
	g := portgroup.New(nil)
	if got := g.Categorize(9999); got != portgroup.CategoryOther {
		t.Fatalf("expected other, got %s", got)
	}
}

func TestCustomOverridesDefault(t *testing.T) {
	custom := map[portgroup.Category][]int{}
	_ = custom
	g := portgroup.New(map[int]portgroup.Category{
		80: portgroup.CategoryOther,
	})
	if got := g.Categorize(80); got != portgroup.CategoryOther {
		t.Fatalf("expected other from custom override, got %s", got)
	}
}

func TestGroupPortsPartitionsCorrectly(t *testing.T) {
	g := portgroup.New(nil)
	ports := []int{80, 443, 3306, 22, 9999}
	groups := g.GroupPorts(ports)

	if len(groups[portgroup.CategoryWeb]) != 2 {
		t.Fatalf("expected 2 web ports, got %d", len(groups[portgroup.CategoryWeb]))
	}
	if len(groups[portgroup.CategoryDatabase]) != 1 {
		t.Fatalf("expected 1 database port, got %d", len(groups[portgroup.CategoryDatabase]))
	}
	if len(groups[portgroup.CategoryRemote]) != 1 {
		t.Fatalf("expected 1 remote port, got %d", len(groups[portgroup.CategoryRemote]))
	}
	if len(groups[portgroup.CategoryOther]) != 1 {
		t.Fatalf("expected 1 other port, got %d", len(groups[portgroup.CategoryOther]))
	}
}

func TestGroupPortsEmptyInput(t *testing.T) {
	g := portgroup.New(nil)
	groups := g.GroupPorts([]int{})
	if len(groups) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(groups))
	}
}

func TestLabelContainsPortAndCategory(t *testing.T) {
	g := portgroup.New(nil)
	label := g.Label(443)
	if label != "443 (web)" {
		t.Fatalf("unexpected label: %s", label)
	}
}
