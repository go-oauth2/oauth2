package store

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/models"
	"golang.org/x/crypto/bcrypt"
)

// Hasher is an interface for hashing and verifying client secrets.
type Hasher interface {
	// Hash hashes the given secret and returns the hashed value.
	Hash(secret string) (string, error)
	// Verify checks if the hashed secret matches the given secret.
	Verify(hashedPassword, secret string) error
}

// BcryptHasher is a Hasher implementation using bcrypt for hashing and verifying secrets.
type BcryptHasher struct{}

func (b *BcryptHasher) Hash(secret string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (b *BcryptHasher) Verify(hashed, secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(secret))
}

// ClientInfoWithHash wraps an oauth2.ClientInfo and provides secret verification using a Hasher.
type ClientInfoWithHash struct {
	wrapped oauth2.ClientInfo
	hasher  Hasher
}

// NewClientInfoWithHash creates a new instance of client info supporting hashed secret verification.
func NewClientInfoWithHash(
	info oauth2.ClientInfo,
	hasher Hasher,
) *ClientInfoWithHash {
	if info == nil {
		return nil
	}
	return &ClientInfoWithHash{
		wrapped: info,
		hasher:  hasher,
	}
}

// VerifyPassword verifies the given plain secret against the hashed secret.
// It implements the oauth2.ClientPasswordVerifier interface.
func (v *ClientInfoWithHash) VerifyPassword(secret string) bool {
	if secret == "" {
		return false
	}
	err := v.hasher.Verify(v.GetSecret(), secret)
	return err == nil
}

// GetID returns the client ID.
func (v *ClientInfoWithHash) GetID() string {
	return v.wrapped.GetID()
}

// GetSecret returns the hashed client secret.
func (v *ClientInfoWithHash) GetSecret() string {
	return v.wrapped.GetSecret()
}

// GetDomain returns the client domain.
func (v *ClientInfoWithHash) GetDomain() string {
	return v.wrapped.GetDomain()
}

// GetUserID returns the user ID associated with the client.
func (v *ClientInfoWithHash) GetUserID() string {
	return v.wrapped.GetUserID()
}

// IsPublic returns true if the client is public.
func (v *ClientInfoWithHash) IsPublic() bool {
	return v.wrapped.IsPublic()
}

// ClientStoreWithHash is a wrapper around oauth2.SavingClientStore that hashes client secrets.
type ClientStoreWithHash struct {
	underlying oauth2.SavingClientStore
	hasher     Hasher
}

// NewClientStoreWithBcrypt creates a new ClientStoreWithHash using bcrypt for hashing.
//
// It is a convenience function for creating a store with the default bcrypt hasher.
// The store will hash client secrets using bcrypt before saving them and would
// return secret information supporting secret verification against the hashed secret.
func NewClientStoreWithBcrypt(store oauth2.SavingClientStore) *ClientStoreWithHash {
	return NewClientStoreWithHash(store, &BcryptHasher{})
}

func NewClientStoreWithHash(underlying oauth2.SavingClientStore, hasher Hasher) *ClientStoreWithHash {
	if hasher == nil {
		hasher = &BcryptHasher{}
	}
	return &ClientStoreWithHash{
		underlying: underlying,
		hasher:     hasher,
	}
}

// GetByID retrieves client information by ID and returns a ClientInfoWithHash instance.
func (w *ClientStoreWithHash) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	info, err := w.underlying.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	rval := NewClientInfoWithHash(info, w.hasher)
	if rval == nil {
		return nil, errors.ErrInvalidClient
	}
	return rval, nil
}

// Save hashes the client secret before saving it to the underlying store.
func (w *ClientStoreWithHash) Save(
	ctx context.Context,
	info oauth2.ClientInfo,
) error {
	if info == nil {
		return errors.ErrInvalidClient
	}
	if info.GetSecret() == "" {
		return errors.ErrInvalidClient
	}

	hashed, err := w.hasher.Hash(info.GetSecret())
	if err != nil {
		return err
	}
	hashedInfo := models.Client{
		ID:     info.GetID(),
		Secret: hashed,
		Domain: info.GetDomain(),
		UserID: info.GetUserID(),
		Public: info.IsPublic(),
	}
	return w.underlying.Save(ctx, &hashedInfo)
}
