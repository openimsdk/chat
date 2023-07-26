package common

import (
	"github.com/OpenIMSDK/tools/errs"
	"regexp"
)

func EmailCheck(email string) error {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	if err := regexMatch(pattern, email); err != nil {
		return errs.Wrap(err, "Email is invalid")
	}
	return nil
}

func AreaCodeCheck(areaCode string) error {
	pattern := `\+[1-9][0-9]{1,2}`
	if err := regexMatch(pattern, areaCode); err != nil {
		return errs.Wrap(err, "AreaCode is invalid")
	}
	return nil
}

func PhoneNumberCheck(phoneNumber string) error {
	pattern := `^1[1-9]{10}`
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
