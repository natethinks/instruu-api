package auth

import "testing"

func TestGeneratePasswordHash(t *testing.T) {
	password := "testing"

	hash := GeneratePasswordHash([]byte(password))

	match := VerifyPassword(hash, []byte(password))
	if !match {
		t.Fail()
	}
	return
}

func TestVerifyPassword(t *testing.T) {
	password := "testing"
	hash := "$2a$04$JouxSY1cV566txEjiDcFOOu.G2H2t8UXAUyzrP8qqZrw.7fsAMEvi"

	match := VerifyPassword(hash, []byte(password))
	if !match {
		t.Fail()
	}
	return
}
