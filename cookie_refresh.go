package railgun

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type authenticationConfig struct {
	SessData     string
	Csrf         string
	Source       string
	RefreshToken string
	RefreshCsrf  string
}

// csrf param is bili_jct of cookie
// default source param is `main_web`
// refresh_token is ac_time_value in local storage
func (r *RailGun) SetAuthentication(sessdata, csrf, source, refreshToken string) error {
	if source == `` {
		source = DEFAULT_REFRESH_COOKIE_SOURCE
	}

	r.authentication = &authenticationConfig{
		SessData:     sessdata,
		Csrf:         csrf,
		Source:       source,
		RefreshToken: refreshToken,
	}

	return r.renewCookie(sessdata)
}

func (r *RailGun) needRenewCookie(sessdata string) (int, error) {
	r.client.setSessData(sessdata)

	resp, err := r.client.checkRefresh()
	if err != nil {
		return 0, err
	}

	if resp.Code != 0 {
		return 0, errors.New(`expired sessdata`)
	}

	if resp.Data.Refresh {
		return resp.Data.Timestamp, nil
	}

	return 0, nil
}

func (r *RailGun) renewCookie(sessdata string) error {
	timestamp, err := r.needRenewCookie(sessdata)
	if err != nil {
		return err
	}

	if timestamp == 0 {
		return nil
	}

	path, err := getCorrespondPath(timestamp)
	if err != nil {
		return err
	}

	err = r.getRefreshCsrf(path)
	if err != nil {
		return err
	}

	return r.refreshCookie()
}

func getCorrespondPath(timestamp int) (string, error) {
	crypted, err := encryptOAEP([]byte(
		stringBuilder(`refresh_`, strconv.Itoa(timestamp)),
	))

	if err != nil {
		return ``, err
	}

	ret := hex.EncodeToString(crypted)

	return ret, nil
}

const RENEW_PUBLIC_CERT = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLgd2OAkcGVtoE3ThUREbio0Eg
Uc/prcajMKXvkCKFCWhJYJcLkcM2DKKcSeFpD/j6Boy538YXnR6VhcuUJOhH2x71
nzPjfdTcqMz7djHum0qSZA0AyCBDABUqCrfNgCiJ00Ra7GmRj+YCK1NJEuewlb40
JNrRuoEUXpabUzGB8QIDAQAB
-----END PUBLIC KEY-----`

func encryptOAEP(password []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(RENEW_PUBLIC_CERT))
	if block == nil {
		err := errors.New(`failed to parse certificate PEM`)
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey := pub.(*rsa.PublicKey)

	h := sha256.New() // sha1.New() or md5.New()
	return rsa.EncryptOAEP(h, rand.Reader, rsaPublicKey, password, nil)
}

func decryptOAEP(cipherdata []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(RENEW_PUBLIC_CERT))
	if block == nil {
		err := errors.New(`failed to parse certificate PEM`)
		return nil, err
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) // ASN.1 PKCS#1 DER encoded form.
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	return rsa.DecryptOAEP(h, rand.Reader, priv, cipherdata, nil)
}

func (r *RailGun) getRefreshCsrf(path string) error {
	regex := regexp.MustCompile(`<div id="1-name">([\w])*</div>`)

	content, err := r.client.getRefreshCsrf(path)
	if err != nil {
		return err
	}

	content = regex.FindString(content)
	if content == `` {
		return errors.New(`failed to get refresh csrf`)
	}

	content = strings.TrimPrefix(content, `<div id="1-name">`)
	content = strings.TrimSuffix(content, `</div>`)
	if content == `` {
		return errors.New(`failed to get refresh csrf`)
	}

	r.authentication.RefreshCsrf = content
	return nil
}

func (a *authenticationConfig) getRefreshCookieParams() map[string]string {
	return map[string]string{
		`csrf`:          a.Csrf,
		`source`:        a.Source,
		`refresh_token`: a.RefreshToken,
		`refresh_csrf`:  a.RefreshCsrf,
	}
}

func (r *RailGun) refreshCookie() error {
	csrf, sessData, refreshToken, err := r.client.refreshCookie(r.authentication.getRefreshCookieParams())
	if err != nil {
		return err
	}

	r.authentication.Csrf = csrf
	r.authentication.SessData = sessData
	r.authentication.RefreshToken = refreshToken
	return r.renewCookie(sessData)
}

func (a *authenticationConfig) Dump() map[string]string {
	return map[string]string{
		`SESSDATA`:      a.SessData,
		`bili_jct`:      a.Csrf,
		`refresh_token`: a.RefreshToken,
	}
}
