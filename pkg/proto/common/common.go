package common

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"regexp"
)

func EmailCheck(email string) error {
	pattern := `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`
	if err := regexMatch(pattern, email); err != nil {
		return errs.Wrap(err, "Email is invalid")
	}
	return nil
}

func AreaCodeCheck(areaCode string) error {
	pattern := ``
	if err := regexMatch(pattern, areaCode); err != nil {
		return errs.Wrap(err, "AreaCode is invalid")
	}
	return nil
}

func PhoneNumberCheck(phoneNumber string) error {
	pattern := ``
	if err := regexMatch(pattern, phoneNumber); err != nil {
		return errs.Wrap(err, "phoneNumber is invalid")
	}
	return nil
}

func regexMatch(pattern string, target string) error {
	reg := regexp.MustCompile(pattern)
	ok := reg.MatchString(target)
	if !ok {
		return errs.ErrArgs
	}
	return nil
}
