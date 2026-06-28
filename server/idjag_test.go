package server

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	oauth2pkg "github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/golang-jwt/jwt/v5"
)

const (
	testIssuer   = "https://idp.example.com"
	testAudience = "https://as.example.com"
	testClientID = "client-abc"
	testSubject  = "user@example.com"
)

func mustGenRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func staticKeyResolver(pub crypto.PublicKey) IssuerKeyResolver {
	return func(_ context.Context, _, _ string) ([]crypto.PublicKey, error) {
		return []crypto.PublicKey{pub}, nil
	}
}

func baseConfig(key *rsa.PrivateKey) *Config {
	return &Config{
		TrustedIDJAGIssuers:    []string{testIssuer},
		IDJAGAudience:          testAudience,
		IDJAGIssuerKeyResolver: staticKeyResolver(&key.PublicKey),
	}
}

func signRS256(t *testing.T, key *rsa.PrivateKey, claims jwt.MapClaims, headerOverrides map[string]any) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["typ"] = IDJAGType
	for k, v := range headerOverrides {
		tok.Header[k] = v
	}
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// noneAlgToken crafts a JWT with alg:none (no valid signature).
func noneAlgToken(claims map[string]any) string {
	h, _ := json.Marshal(map[string]any{"alg": "none", "typ": IDJAGType})
	c, _ := json.Marshal(claims)
	return base64.RawURLEncoding.EncodeToString(h) + "." +
		base64.RawURLEncoding.EncodeToString(c) + "."
}

func validClaims(jti string, expOffset time.Duration) jwt.MapClaims {
	now := time.Now()
	return jwt.MapClaims{
		"iss":       testIssuer,
		"sub":       testSubject,
		"aud":       []string{testAudience},
		"client_id": testClientID,
		"jti":       jti,
		"exp":       now.Add(expOffset).Unix(),
		"iat":       now.Unix(),
	}
}

func TestValidateIDJAGAssertion(t *testing.T) {
	key := mustGenRSAKey(t)
	cfg := baseConfig(key)

	tests := []struct {
		name      string
		assertion func() string
		clientID  string
		cfg       *Config
		wantErr   error
	}{
		{
			name:      "valid RS256 assertion",
			assertion: func() string { return signRS256(t, key, validClaims("jti-ok", 5*time.Minute), nil) },
			clientID:  testClientID,
			cfg:       cfg,
		},
		{
			name: "typ JWT rejected",
			assertion: func() string {
				return signRS256(t, key, validClaims("jti-typ-jwt", 5*time.Minute), map[string]any{"typ": "JWT"})
			},
			clientID: testClientID,
			cfg:      cfg,
			wantErr:  errors.ErrInvalidGrant,
		},
		{
			name: "typ id_token rejected",
			assertion: func() string {
				return signRS256(t, key, validClaims("jti-typ-idtoken", 5*time.Minute), map[string]any{"typ": "id_token"})
			},
			clientID: testClientID,
			cfg:      cfg,
			wantErr:  errors.ErrInvalidGrant,
		},
		{
			name: "alg none rejected",
			assertion: func() string {
				now := time.Now()
				return noneAlgToken(map[string]any{
					"iss": testIssuer, "sub": testSubject, "aud": []string{testAudience},
					"client_id": testClientID, "jti": "jti-none",
					"exp": now.Add(5 * time.Minute).Unix(), "iat": now.Unix(),
				})
			},
			clientID: testClientID,
			cfg:      cfg,
			wantErr:  errors.ErrInvalidGrant,
		},
		{
			name: "alg HS256 rejected",
			assertion: func() string {
				tok := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims("jti-hs256", 5*time.Minute))
				tok.Header["typ"] = IDJAGType
				s, err := tok.SignedString([]byte("secret"))
				if err != nil {
					t.Fatal(err)
				}
				return s
			},
			clientID: testClientID,
			cfg:      cfg,
			wantErr:  errors.ErrInvalidGrant,
		},
		{
			name: "untrusted issuer rejected",
			assertion: func() string {
				c := validClaims("jti-bad-iss", 5*time.Minute)
				c["iss"] = "https://evil.example.com"
				return signRS256(t, key, c, nil)
			},
			clientID: testClientID,
			cfg:      cfg,
			wantErr:  errors.ErrInvalidGrant,
		},
		{
			name: "wrong aud rejected",
			assertion: func() string {
				c := validClaims("jti-bad-aud", 5*time.Minute)
				c["aud"] = []string{"https://wrong.example.com"}
				return signRS256(t, key, c, nil)
			},
			clientID: testClientID,
			cfg:      cfg,
			wantErr:  errors.ErrInvalidGrant,
		},
		{
			name:      "client_id mismatch rejected",
			assertion: func() string { return signRS256(t, key, validClaims("jti-cid-mismatch", 5*time.Minute), nil) },
			clientID:  "other-client",
			cfg:       cfg,
			wantErr:   errors.ErrInvalidGrant,
		},
		{
			name:      "expired outside clock skew rejected",
			assertion: func() string { return signRS256(t, key, validClaims("jti-expired", -90*time.Second), nil) },
			clientID:  testClientID,
			cfg:       cfg,
			wantErr:   errors.ErrInvalidGrant,
		},
		{
			// Token expired 30s ago; default clock skew is 60s — should succeed.
			name:      "expired within clock skew accepted",
			assertion: func() string { return signRS256(t, key, validClaims("jti-skew", -30*time.Second), nil) },
			clientID:  testClientID,
			cfg:       cfg,
		},
		{
			name: "missing jti rejected",
			assertion: func() string {
				c := validClaims("", 5*time.Minute)
				delete(c, "jti")
				return signRS256(t, key, c, nil)
			},
			clientID: testClientID,
			cfg:      cfg,
			wantErr:  errors.ErrInvalidGrant,
		},
		{
			// Store is pre-populated with the same jti before validation runs.
			name:      "replay rejected",
			assertion: func() string { return signRS256(t, key, validClaims("jti-replay", 5*time.Minute), nil) },
			clientID:  testClientID,
			cfg: func() *Config {
				c := *cfg
				c.IDJAGReplayStore = oauth2pkg.NewMemoryAssertionReplayStore()
				_ = c.IDJAGReplayStore.StoreAssertionID(context.Background(), "jti-replay", time.Now().Add(5*time.Minute))
				return &c
			}(),
			wantErr: errors.ErrInvalidGrant,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateIDJAGAssertion(context.Background(), tt.cfg, tt.clientID, tt.assertion())
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got %v, want %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateIDJAGAssertionReturnedClaims(t *testing.T) {
	key := mustGenRSAKey(t)
	cfg := baseConfig(key)

	claims := validClaims("jti-claims-check", 5*time.Minute)
	claims["scope"] = "read write"
	claims["resource"] = []string{"https://api.example.com"}

	got, err := validateIDJAGAssertion(context.Background(), cfg, testClientID, signRS256(t, key, claims, nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Issuer != testIssuer {
		t.Errorf("Issuer: got %q want %q", got.Issuer, testIssuer)
	}
	if got.Subject != testSubject {
		t.Errorf("Subject: got %q want %q", got.Subject, testSubject)
	}
	if got.ClientID != testClientID {
		t.Errorf("ClientID: got %q want %q", got.ClientID, testClientID)
	}
	if got.JTI != "jti-claims-check" {
		t.Errorf("JTI: got %q want %q", got.JTI, "jti-claims-check")
	}
	if got.Scope != "read write" {
		t.Errorf("Scope: got %q want %q", got.Scope, "read write")
	}
	if len(got.Resource) != 1 || got.Resource[0] != "https://api.example.com" {
		t.Errorf("Resource: got %v", got.Resource)
	}
}

func TestNewOIDCIssuerKeyResolver(t *testing.T) {
	key := mustGenRSAKey(t)
	kid := "key-1"

	n := base64.RawURLEncoding.EncodeToString(key.PublicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(new(big.Int).SetInt64(int64(key.PublicKey.E)).Bytes())

	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			json.NewEncoder(w).Encode(map[string]string{"jwks_uri": ts.URL + "/jwks"})
		case "/jwks":
			json.NewEncoder(w).Encode(map[string]any{
				"keys": []map[string]any{
					{"kty": "RSA", "kid": kid, "n": n, "e": e},
				},
			})
		}
	}))
	defer ts.Close()

	resolver := NewOIDCIssuerKeyResolver(nil)

	t.Run("matching kid returns key", func(t *testing.T) {
		keys, err := resolver(context.Background(), ts.URL, kid)
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(keys))
		}
		rsaKey, ok := keys[0].(*rsa.PublicKey)
		if !ok {
			t.Fatal("expected *rsa.PublicKey")
		}
		if rsaKey.N.Cmp(key.PublicKey.N) != 0 {
			t.Error("key modulus mismatch")
		}
	})

	t.Run("empty kid returns all keys", func(t *testing.T) {
		keys, err := resolver(context.Background(), ts.URL, "")
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(keys))
		}
	})

	t.Run("non-matching kid returns empty slice", func(t *testing.T) {
		keys, err := resolver(context.Background(), ts.URL, "other-kid")
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != 0 {
			t.Errorf("expected 0 keys, got %d", len(keys))
		}
	})

	t.Run("missing jwks_uri returns error", func(t *testing.T) {
		badTs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{})
		}))
		defer badTs.Close()
		_, err := resolver(context.Background(), badTs.URL, "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("cancelled context returns error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := resolver(ctx, ts.URL, kid)
		if err == nil {
			t.Error("expected error from cancelled context, got nil")
		}
	})
}
