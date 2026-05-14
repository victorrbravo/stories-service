package helpers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/twinj/uuid"

	"stories/constants"
	"stories/types"

	"github.com/dgrijalva/jwt-go"
)

type AccessDetails struct {
	AccessUuid string
	Email      string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const charsetCapital = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type parameters struct {
	Feature string
}

// GUID UUID  helpers functions
func ConvertUUIDToString(uid []byte) string {
	return uuid.New(uid).String()
}

func ConvertUUIDToTypeString(uid []byte) types.String {
	return types.NewString(uuid.New(uid).String())
}

func ConvertJsonToDbArray(values []string) string {
	var str string = "{"
	for _, value := range values {
		str = str + fmt.Sprintf("'%s', ", value)
	}
	str = strings.TrimRight(str, ", ")
	str = str + "}"
	log.Println("Stories", "str", str)
	return str
}

func ConvertDbToJsonArray(value string) []string {
	str := value
	var values []string
	str = strings.TrimLeft(str, "{")
	str = strings.TrimRight(str, "}")
	str = strings.Replace(str, "'", "", -1)
	str = strings.Replace(str, "\"", "", -1)
	values = strings.Split(str, ",")
	return values
}

func ConvertDbToJson(value string) string {
	str := value
	str = strings.Replace(str, "{", "[", -1)
	str = strings.Replace(str, "}", "]", -1)
	return str
}

// isEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func StringWithCharset(length int, charset string) string {
	var seededRand *rand.Rand
	var b []byte
	rand.Seed(time.Now().Unix())
	seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b = make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func CreateJWTToken(email string) (*TokenDetails, error) {
	var err error
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 72).Unix()
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Minute * 24 * 7).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + email
	// Creating Access Token
	value := os.Getenv(constants.JWTAccessEnvironmentOSName)
	if len(value) == 0 {
		_ = os.Setenv(constants.JWTAccessEnvironmentOSName, constants.DefaultJWT) // this should be in an env file
	}
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["email"] = email
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv(constants.JWTAccessEnvironmentOSName)))
	if err != nil {
		return td, err
	}
	// Creating Refresh Token
	_ = os.Setenv(constants.JWTRefreshEnvironmentOSName, constants.DefaultRefreshJWT) // this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["email"] = email
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv(constants.JWTRefreshEnvironmentOSName)))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func checkOSJWTEnvironment() {
	value := os.Getenv(constants.JWTAccessEnvironmentOSName)
	if len(value) == 0 {
		log.Println("JWTAccessEnvironmentOSName is not set, setting it to default value")
		_ = os.Setenv(constants.JWTAccessEnvironmentOSName, constants.DefaultJWT) // this should be in an env file
	}
}

func CreateInvitationJWTToken(email string) (*TokenDetails, error) {
	var err error
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 60 * 24 * 365 * 2).Unix()
	td.AccessUuid = uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Minute * 60 * 24 * 365 * 2).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + email
	// Creating Access Token
	checkOSJWTEnvironment()
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["email"] = email
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv(constants.JWTAccessEnvironmentOSName)))
	if err != nil {
		return td, err
	}
	// Creating Refresh Token
	_ = os.Setenv(constants.JWTRefreshEnvironmentOSName, constants.DefaultRefreshJWT) // this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["email"] = email
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv(constants.JWTRefreshEnvironmentOSName)))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func verifyToken(t string) (*jwt.Token, error) {
	tokenString := t
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(constants.JWTAccessEnvironmentOSName)), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(t string) (jwt.MapClaims, error) {
	checkOSJWTEnvironment()
	token, err := verifyToken(t)
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}

func VerifyRefreshToken(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(constants.JWTRefreshEnvironmentOSName)), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func GenerateCode() string {
	var code string
	sha := sha256.New()
	code = StringWithCharset(rand.Intn(64), charset)
	_, _ = io.WriteString(sha, code)
	return hex.EncodeToString(sha.Sum(nil))
}

func AdjustDateToLocal(current time.Time) (time.Time, error) {
	var err error
	var loc *time.Location
	loc, err = time.LoadLocation("America/Ny")
	if err != nil {
		log.Println("error convert time in GetChallengeDate err", err)
		return current, err
	}
	return current.In(loc), nil
}

func GetSqlFromTemplate(base string, feature string) (string, error) {
	var tpl *template.Template
	var buffer bytes.Buffer
	var err error
	tpl = template.New("update-queries")
	// Try to parse the code
	tpl, err = tpl.Parse(base)
	if err != nil {
		return "", nil
	}
	// WARNING: make sure that non of these parameters are user input
	//          otherwise this it would be possible to perform SQL injection
	//          within this code
	err = tpl.Execute(&buffer, parameters{
		Feature: feature,
	})
	return buffer.String(), nil
}

func PenultimateElement(str string, sep string) (string, error) {
	var ss []string
	var value string
	var size int
	ss = strings.Split(str, sep)
	size = len(ss)
	if size < 2 {
		return value, fmt.Errorf(constants.ArrayMustHave2elementsOrMore, str)
	}
	return ss[size-2], nil
}
