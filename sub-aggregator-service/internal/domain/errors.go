package domain

import "errors"

var ErrAlreadyExists = errors.New("creating object: already exists")
var ErrNotFound = errors.New("subscription not found")
