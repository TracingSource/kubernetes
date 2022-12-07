package renewal

import (
	"crypto/x509"
	"math"
	"testing"
	"time"
)

func TestExpirationInfo(t *testing.T) {
	validity := 365 * 24 * time.Hour
	cert := &x509.Certificate{
		NotAfter: time.Now().Add(validity),
	}

	e := newExpirationInfo("x", cert, false)

	if math.Abs(float64(validity-e.ResidualTime())) > float64(5*time.Second) { // using 5s of tolerance because the function is not deterministic (it uses time.Now()) and we want to avoid flakes
		t.Errorf("expected IsInRenewalWindow equal to %v, saw %v", validity, e.ResidualTime())
	}
}
