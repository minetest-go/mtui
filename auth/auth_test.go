package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test:enter from the database column
const CREDENTIALS = "#1#TxqLUa/uEJvZzPc3A0xwpA#oalXnktlS0bskc7bccsoVTeGwgAwUOyYhhceBu7wAyITkYjCtrzcDg6W5Co5V+oWUSG13y7TIoEfIg6rafaKzAbwRUC9RVGCeYRIUaa0hgEkIe9VkDmpeQ/kfF8zT8p7prOcpyrjWIJR+gmlD8Bf1mrxoPoBLDbvmxkcet327kQ9H4EMlIlv+w3XCufoPGFQ1UrfWiVqqK8dEmt/ldLPfxiK1Rg8MkwswEekymP1jyN9Cpq3w8spVVcjsxsAzI5M7QhSyqMMrIThdgBsUqMBOCULdV+jbRBBiA/ClywtZ8vvBpN9VGqsQuhmQG0h5x3fqPyR2XNdp9Ocm3zHBoJy/w"

func TestCreateDBPasswordAuth(t *testing.T) {
	salt, verifier, err := CreateAuth("x", "y")
	assert.NoError(t, err)
	data := CreateDBPassword(salt, verifier)
	salt, verifier, err = ParseDBPassword(data)
	assert.NoError(t, err)
	ok, err := VerifyAuth("x", "y", salt, verifier)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestInvalidVersion(t *testing.T) {
	_, _, err := ParseDBPassword("#2#123#456")
	assert.Error(t, err)
}

func TestInvalidDelimiter(t *testing.T) {
	_, _, err := ParseDBPassword("#1#123")
	assert.Error(t, err)
}

func TestInvalidBase64(t *testing.T) {
	_, _, err := ParseDBPassword("#1#...#...")
	assert.Error(t, err)
	_, _, err = ParseDBPassword("#1#TxqLUa/uEJvZzPc3A0xwpA#...")
	assert.Error(t, err)
}

func TestAuth(t *testing.T) {
	salt, verifier, err := ParseDBPassword(CREDENTIALS)
	assert.NoError(t, err)

	// valid credentials
	success, err := VerifyAuth("test", "enter", salt, verifier)
	assert.NoError(t, err)
	assert.True(t, success)

	// invalid password
	success, err = VerifyAuth("test", "bogus", salt, verifier)
	assert.NoError(t, err)
	assert.False(t, success)

	// invalid user
	success, err = VerifyAuth("testx", "enter", salt, verifier)
	assert.NoError(t, err)
	assert.False(t, success)

	// ad-hoc creation
	salt, verifier, err = CreateAuth("a", "b")
	assert.NoError(t, err)

	// valid credentials
	success, err = VerifyAuth("a", "b", salt, verifier)
	assert.NoError(t, err)
	assert.True(t, success)
}

func TestLagacyAuth(t *testing.T) {
	success := VerifyLegacyPassword("BuckarooBanzai", "enter", "vTELhF7/I72pH9rsz79yU/1hZ5A")
	assert.True(t, success)

	success = VerifyLegacyPassword("BuckarooBanzai", "blah", "vTELhF7/I72pH9rsz79yU/1hZ5A")
	assert.False(t, success)
}

func TestValidateUsername(t *testing.T) {
	valid := []string{
		"my_user01",
		"somedude",
		"_",
		"-",
		"01234567890123456789",
	}
	invalid := []string{
		"",
		" ",
		"012345678901234567890",
		"my name",
		" x",
		"x ",
		"%",
		"<",
		"*",
		".",
		"|",
	}

	for _, str := range valid {
		assert.NoError(t, ValidateUsername(str), "string was not valid: '%s'", str)
	}
	for _, str := range invalid {
		assert.Error(t, ValidateUsername(str), "string was not invalid: '%s'", str)
	}
}
