package portcorrelate

import (
	"testing"
	"time"
)

func TestObserveCreatesPair(t *testing.T) {
	c := New(time.Second)
	c.Observe(80, 443)
	pairs := c.All()
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if pairs[0].A != 80 || pairs[0].B != 443 {
		t.Errorf("unexpected pair: %+v", pairs[0])
	}
}

func TestObserveNormalisesOrder(t *testing.T) {
	c := New(time.Second)
	c.Observe(443, 80)
	pairs := c.All()
	if pairs[0].A != 80 || pairs[0].B != 443 {
		t.Errorf("expected A<B, got %+v", pairs[0])
	}
}

func TestObserveIncrementsCount(t *testing.T) {
	c := New(time.Second)
	c.Observe(80, 443)
	c.Observe(80, 443)
	c.Observe(443, 80)
	pairs := c.All()
	if pairs[0].Count != 3 {
		t.Errorf("expected count 3, got %d", pairs[0].Count)
	}
}

func TestSamePairIgnored(t *testing.T) {
	c := New(time.Second)
	c.Observe(80, 80)
	if len(c.All()) != 0 {
		t.Error("expected no pair for identical ports")
	}
}

func TestStrongFiltersLowCount(t *testing.T) {
	c := New(time.Second)
	c.Observe(80, 443)
	c.Observe(22, 8080)
	c.Observe(22, 8080)
	strong := c.Strong(2)
	if len(strong) != 1 {
		t.Fatalf("expected 1 strong pair, got %d", len(strong))
	}
	if strong[0].A != 22 || strong[0].B != 8080 {
		t.Errorf("unexpected strong pair: %+v", strong[0])
	}
}

func TestResetClearsAll(t *testing.T) {
	c := New(time.Second)
	c.Observe(80, 443)
	c.Reset()
	if len(c.All()) != 0 {
		t.Error("expected empty after reset")
	}
}

func TestMultiplePairsAreIndependent(t *testing.T) {
	c := New(time.Second)
	c.Observe(80, 443)
	c.Observe(22, 8080)
	if len(c.All()) != 2 {
		t.Errorf("expected 2 pairs, got %d", len(c.All()))
	}
}
