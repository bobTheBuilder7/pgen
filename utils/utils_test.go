package utils

import (
	"testing"

	"github.com/bobTheBuilder7/assert"
)

func TestToPascalCase_SingleWord(t *testing.T) {
	assert.Equal(t, ToPascalCase("name"), "Name")
}

func TestToPascalCase_TwoWords(t *testing.T) {
	assert.Equal(t, ToPascalCase("user_id"), "UserId")
}

func TestToPascalCase_ThreeWords(t *testing.T) {
	assert.Equal(t, ToPascalCase("created_at_utc"), "CreatedAtUtc")
}

func TestToPascalCase_AlreadyPascal(t *testing.T) {
	assert.Equal(t, ToPascalCase("Name"), "Name")
}

func TestToPascalCase_AllUppercase(t *testing.T) {
	assert.Equal(t, ToPascalCase("ID"), "Id")
}

func TestToPascalCase_LeadingUnderscore(t *testing.T) {
	assert.Equal(t, ToPascalCase("_name"), "Name")
}

func TestToPascalCase_TrailingUnderscore(t *testing.T) {
	assert.Equal(t, ToPascalCase("name_"), "Name")
}

func TestToPascalCase_Empty(t *testing.T) {
	assert.Equal(t, ToPascalCase(""), "")
}

func TestToPascalCase_ConsecutiveUnderscores(t *testing.T) {
	assert.Equal(t, ToPascalCase("user__name"), "UserName")
}

func TestToPascalCase_MixedCase(t *testing.T) {
	assert.Equal(t, ToPascalCase("userNAME_AGE"), "UsernameAge")
}

func TestToPascalCase_SingleChar(t *testing.T) {
	assert.Equal(t, ToPascalCase("a_b_c"), "ABC")
}

func TestToPascalCase_NumbersInWords(t *testing.T) {
	assert.Equal(t, ToPascalCase("int8_value"), "Int8Value")
}

func TestToPascalCase_OnlyUnderscores(t *testing.T) {
	assert.Equal(t, ToPascalCase("___"), "")
}
