package auth

import (
	"time"

	"github.com/tinode/chat/server/store/types"
)

const (
	// No error
	NoErr = iota
	// No change
	InfoNotModified
	// DB or other internal failure
	ErrInternal
	// The secret cannot be parsed or otherwise wrong
	ErrMalformed
	// Authentication failed (wrong password)
	ErrFailed
	// Duplicate credential
	ErrDuplicate
	// The operation is unsupported
	ErrUnsupported
	// Secret has expired
	ErrExpired
	// Policy violation, e.g. password too weak.
	ErrPolicy
)

// Authentication levels
const (
	// Undefined/not authenticated
	LevelNone = iota * 10
	// Anonymous user/light authentication
	LevelAnon
	// Fully authenticated user
	LevelAuth
	// Superuser (currently unused)
	LevelRoot
)

// Interface which auth providers must implement
type AuthHandler interface {
	// Initialize the handler
	Init(jsonconf string) error

	// Add persistent record to database. Returns a numeric error code to indicate
	// if the error is due to a duplicate or some other error.
	// store.AddAuthRecord("scheme", "unique", "secret")
	// Returns: auth level, error code, error
	AddRecord(uid types.Uid, secret []byte, lifetime time.Duration) (int, int, error)

	// Update existing record with new credentials. Returns a numeric error code to indicate
	// if the error is due to a duplicate or some other error.
	// store.UpdateAuthRecord("scheme", "unique", "secret")
	UpdateRecord(uid types.Uid, secret []byte, lifetime time.Duration) (int, error)

	// Given a user-provided authentication secret (such as "login:password"
	// return user ID, time when the secret expires (zero, if never) or an error code.
	// store.Users.GetAuthRecord("scheme", "unique")
	// Returns: user ID, user auth level, token expiration time, error code.
	Authenticate(secret []byte) (types.Uid, int, time.Time, int)

	// Verify if the provided secret can be considered unique by the auth scheme
	// E.g. if login is unique.
	// store.GetAuthRecord(scheme, unique)
	IsUnique(secret []byte) (bool, error)

	// Generate a new secret, if appropriate.
	GenSecret(uid types.Uid, authLvl int, lifetime time.Duration) ([]byte, time.Time, int)
}

func AuthLevelName(authLvl int) string {
	switch authLvl {
	case LevelNone:
		return "none"
	case LevelAnon:
		return "anon"
	case LevelAuth:
		return "auth"
	case LevelRoot:
		return "root"
	default:
		return "unkn"
	}
}
