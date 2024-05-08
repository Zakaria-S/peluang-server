package domain

import "errors"

//auth error
var (
	ErrEmailExist        = errors.New("email already exists")
	ErrInvalidOTP        = errors.New("invalid otp")
	ErrExpiredOTP        = errors.New("otp has expired")
	ErrInvalidCredential = errors.New("invalid credential")
)

//token error
var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrWrongTypeToken = errors.New("wrong type token")
	ErrExpiredToken   = errors.New("token has expired")
	ErrEmptyToken     = errors.New("empty token")
	ErrNoBerearToken  = errors.New("token no Bearer")
)

var (
	ErrMinMax = errors.New("this price range already exists")
)
