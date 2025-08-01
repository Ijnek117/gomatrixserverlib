package spec

import (
	"fmt"
	"regexp"
	"strings"
)

const UserSigil = '@'
const LocalDomainSeparator = ':'

var validEncryptedUsernameRegex = regexp.MustCompile(`^[0-9a-zA-Z+/]+={0,2}$`)

// A UserID identifies a matrix user as per the matrix specification
type EncryptedUserID struct {
	raw    string
	local  string
	domain string
}

// Creates a new Encrypted UserID, returning an error if invalid
func NewEncryptedUserID(id string, allowHistoricalIDs bool) (*EncryptedUserID, error) {
	// Do we make it adhere to current UserID requirements?
	// Will we need to base64 encode it?
	return parseAndValidateEncryptedUserID(id, allowHistoricalIDs)
}

// Creates a new UserID, panicing if invalid
func NewEncryptedUserIDOrPanic(id string, allowHistoricalIDs bool) EncryptedUserID {
	encryptedUserID, err := parseAndValidateEncryptedUserID(id, allowHistoricalIDs)
	if err != nil {
		panic(fmt.Sprintf("NewUserIDOrPanic failed: invalid encrypted user ID %s: %s", id, err.Error()))
	}
	return *encryptedUserID
}

// Returns the full userID string including leading sigil
func (user *EncryptedUserID) String() string {
	return user.raw
}

// Returns just the localpart of the userID
func (user *EncryptedUserID) Local() string {
	return user.local
}

// Returns just the domain of the userID
func (user *EncryptedUserID) Domain() ServerName {
	return ServerName(user.domain)
}

func parseAndValidateEncryptedUserID(id string, allowHistoricalIDs bool) (*EncryptedUserID, error) {
	idLength := len(id)
	// FIXME: 750 chosen arbitrarly for now
	if idLength < 4 || idLength > 750 { // 4 since minimum userID includes an @, :, non-empty localpart, non-empty domain
		return nil, fmt.Errorf("length %d is not within the bounds 4-750", idLength)
	}
	if id[0] != userSigil {
		return nil, fmt.Errorf("first character is not '%c'", userSigil)
	}

	localpart, domain, found := strings.Cut(id[1:], string(localDomainSeparator))
	if !found {
		return nil, fmt.Errorf("at least one '%c' is expected in the user id", localDomainSeparator)
	}
	if _, _, ok := ParseAndValidateServerName(ServerName(domain)); !ok {
		return nil, fmt.Errorf("domain is invalid")
	}

	if allowHistoricalIDs {
		// NOTE: Allowed historical userIDs:
		// https://spec.matrix.org/v1.4/appendices/#historical-user-ids
		if !historicallyValidCharacters(localpart) {
			return nil, fmt.Errorf("local part contains invalid characters from historical set")
		}
	} else {
		//base 64 decoding

		// NOTE: Allowed in the latest spec:
		// https://spec.matrix.org/v1.4/appendices/#user-identifiers
		if !validEncryptedUsernameRegex.MatchString(localpart) {
			fmt.Printf("\nThis is the extracted localpart: %v\n", localpart)
			return nil, fmt.Errorf("local part contains invalid characters")
		}
	}

	encryptedUserID := &EncryptedUserID{
		raw:    id,
		local:  localpart,
		domain: domain,
	}
	return encryptedUserID, nil
}
