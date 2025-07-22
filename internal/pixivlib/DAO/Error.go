package DAO

import "errors"

var NotFound = errors.New("Illust Not Found")

type MyError interface {
	error
	GetS() string
}

type NotGood struct {
	S   string
	Err error
}

func (i *NotGood) Error() string {
	return i.S
}

func (i *NotGood) Unwrap() error {
	return i.Err
}

func (i *NotGood) GetS() string {
	return i.S
}

type AgeLimit struct {
	S   string
	Err error
}

func (i *AgeLimit) Error() string {
	return i.S
}

func (i *AgeLimit) Unwrap() error {
	return i.Err
}

func (i *AgeLimit) GetS() string {
	return i.S
}

type TooFastRequest struct {
	S   string
	Err error
}

func (i *TooFastRequest) Error() string {
	return i.S
}

func (i *TooFastRequest) Unwrap() error {
	return i.Err
}

func ContainMyerror(err error) bool {
	var check MyError
	check = new(AgeLimit)
	return errors.As(err, &check)
}
