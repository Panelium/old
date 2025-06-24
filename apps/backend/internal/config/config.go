package config

import "time"

const Pepper string = "s0m3R4nd0mP3pp3r"     // needs to be changed to crypto/rand.Text()
const JWTSecret string = "0th3rR4nd0mS3cr3t" // needs to be changed to actual RSA key

const AccessTokenDuration = 5 * time.Minute // 5 minutes
const RefreshTokenDuration = 24 * time.Hour // 24 hours
