package jptel

import (
	"errors"
	"strings"
)

// from: http://www.soumu.go.jp/main_sosiki/joho_tsusin/top/tel_number/number_shitei.html
//  0 + [xxxx]  ~  [xxxx] + [xxxx]
//  1 + <- total 5 num -> = 6

const (
	totalCodeLen = 6

	mobileCodePrefixLen = 3
	mobileCodeLen       = 4
	freeDialPrefixLen   = 4
	freeDialCodeLen     = 3
	otherCodePrefixLen  = 4
	otherCodeLen        = 3
)

var mobileCode = []string{
	"050", // IP電話
	"070", // 携帯電話/PHS
	"080", // 携帯電話
	"090", // 携帯電話
	"020", // その他
}

var freeDialCode = []string{
	"0120", // フリーダイヤル 0120-xxx-xxx
	"0800", // フリーダイヤル 0800-xxx-xxxx
}

var otherCode = []string{
	"0570", // ナビダイヤル 0570-xxx-xxx
	"0990", // ダイヤルQ2 0990-xxx-xxx
}

var areaCodes = [][]string{
	areaCode5,
	areaCode4,
	areaCode3,
	areaCode2,
}

var (
	shortError = errors.New("telephone number is too short.")
	matchError = errors.New("telephone number does not match any area code.")
)

// Split splits japaneses telephone number to a slice consist of AreaCode, CityCode, SubscriberNumber.
func Split(tel string) (JPTelephoneNumber, error) {
	// 固定電話
	for _, areaCode := range areaCodes {
		for _, code := range areaCode {
			if strings.HasPrefix(tel, code) {
				if len(tel) < totalCodeLen {
					return JPTelephoneNumber{}, shortError
				}
				codeLen := len(code)
				return JPTelephoneNumber{
					AreaCode:       tel[:codeLen],
					CityCode:       tel[codeLen:totalCodeLen],
					SubscriberCode: tel[totalCodeLen:],
				}, nil
			}
		}
	}

	// フリーダイヤル
	for _, code := range freeDialCode {
		if strings.HasPrefix(tel, code) {
			if len(tel) < freeDialPrefixLen+freeDialCodeLen {
				return JPTelephoneNumber{}, shortError
			}
			return JPTelephoneNumber{
				AreaCode:       tel[:freeDialPrefixLen],
				CityCode:       tel[freeDialPrefixLen : freeDialPrefixLen+freeDialCodeLen],
				SubscriberCode: tel[freeDialPrefixLen+freeDialCodeLen:],
			}, nil
		}
	}

	// 携帯番号
	for _, code := range mobileCode {
		if strings.HasPrefix(tel, code) {
			if len(tel) < mobileCodePrefixLen+mobileCodeLen {
				return JPTelephoneNumber{}, shortError
			}
			return JPTelephoneNumber{
				AreaCode:       tel[:mobileCodePrefixLen],
				CityCode:       tel[mobileCodePrefixLen : mobileCodePrefixLen+mobileCodeLen],
				SubscriberCode: tel[mobileCodePrefixLen+mobileCodeLen:],
			}, nil
		}
	}

	// その他番号
	for _, code := range otherCode {
		if strings.HasPrefix(tel, code) {
			if len(tel) < otherCodePrefixLen+otherCodeLen {
				return JPTelephoneNumber{}, shortError
			}
			return JPTelephoneNumber{
				AreaCode:       tel[:otherCodePrefixLen],
				CityCode:       tel[otherCodePrefixLen : otherCodePrefixLen+otherCodeLen],
				SubscriberCode: tel[otherCodePrefixLen+otherCodeLen:],
			}, nil
		}
	}

	return JPTelephoneNumber{}, matchError
}