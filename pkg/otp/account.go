package otp

import (
	"github.com/pquerna/otp"
	"github.com/shumin1027/otpd/pkg/badger"
	"github.com/shumin1027/otpd/pkg/logger"
	"github.com/vmihailenco/msgpack/v5"
)

var bucket *badger.Bucket

func Init(path string) {
	stor, _ := badger.Open(path, logger.L())
	bucket = stor.CreateBucket("otp")
}

type Account struct {
	OTP    string `json:"otp"`
	Name   string `json:"name"`
	QRCode string `json:"qr_code"`
}

func (account *Account) Key() (*otp.Key, error) {
	key, err := otp.NewKeyFromURL(account.OTP)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (account *Account) Save() error {
	// encode
	buf, err := msgpack.Marshal(&account)
	if err != nil {
		return err
	}
	return bucket.Set([]byte(account.Name), buf)
}

func Get(name string) (*Account, error) {
	if bucket.Has([]byte(name)) {
		val, err := bucket.Get([]byte(name))
		if err != nil {
			return nil, err
		}
		var account Account
		// decode
		err = msgpack.Unmarshal(val, &account)
		if err != nil {
			return nil, err
		}
		return &account, nil
	}
	return nil, nil
}
