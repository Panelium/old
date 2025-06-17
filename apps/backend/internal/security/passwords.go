package security

func HashPassword(password string, salt string, pepper string) (string, error) {
	return "", nil
}

func VerifyPassword(hashedPassword string, password string, salt string, pepper string) (bool, error) {
	return false, nil
}

func GenerateSalt() (string, error) {
	return "", nil
}

func GeneratePepper() (string, error) {
	return "", nil
}
