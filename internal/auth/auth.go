package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

//go:generate mockery --name IAuthService
type IAuthService interface {
	Authorize(r *http.Request) (user string, err error)
}

type authService struct {

}

var _ IAuthService = (*authService)(nil)

func NewAuthService() IAuthService {
	return &authService{}
}

var (
	stringType         = reflect.TypeOf("")
	interfaceSliceType = reflect.TypeOf([]interface{}{})
	authHeaderName     = "authorization"
	bearerSchema       = "Bearer "
	internalSessionId  = "SLB_SESSION_ID"
)

// Claims struct
// Issuer added for user role service
type SlbCustomClaims struct {
	AZP         string      `json:"azp,omitempty"`
	RawAudClaim interface{} `json:"aud,omitempty"`
	Email       string      `json:"email,omitempty"`
	UserID      string      `json:"userid,omitempty"`
	JTI         string      `json:"jti,omitempty"`
	ApiKey      string      `json:"apiKey,omitempty"`
	Sub         string      `json:"sub,omitempty"`
	Issuer      string      `json:"iss,omitempty"`
	Aud         string
	Audiences   []string
}

func (svc *authService) Authorize(r *http.Request) (user string, err error) {
	urlPath := r.URL.Path
	// differentiate if a request is called from external network or internal network
	isInternalRequest := urlPath == "/"
	if !isInternalRequest {
		user, err := svc.getUserFromHttpRequest(r)
		return user, err
	} else {
		// Use SLB_SESSION_ID to identify the client from internal caller
		sessionId := r.Header.Get(internalSessionId)
		if sessionId == "" {
			return "", NewUnauthorizedError("missing header: SLB_SESSION_ID")
		} else {
			return sessionId, nil
		}
	}
}

func (svc *authService) getUserFromHttpRequest(r *http.Request) (user string, err error) {
	authHeaderValue := r.Header.Get(authHeaderName)
	var bearerToken string
	if len(authHeaderValue) > len(bearerSchema) {
		bearerToken = authHeaderValue[len(bearerSchema):]
	}
	if bearerToken == "" {
		return "", NewUnauthorizedError("empty bearerToken")
	}
	tokenBody := strings.Split(bearerToken, `.`)
	if len(tokenBody) < 2 {
		return "", NewUnauthorizedError("unexpected number of parts")
	}
	data, err := base64.RawURLEncoding.DecodeString(tokenBody[1])
	if err != nil {
		return "", NewUnauthorizedError(err.Error())
	}
	claims := &SlbCustomClaims{}
	err = json.Unmarshal(data, claims)
	if err != nil {
		return "", NewUnauthorizedError(err.Error())
	}
	if claims.RawAudClaim != nil {
		// safety check before type assertion
		audClaimType := reflect.TypeOf(claims.RawAudClaim)
		if audClaimType == stringType {
			claims.Aud = claims.RawAudClaim.(string)
		} else if audClaimType == interfaceSliceType {
			for _, a := range claims.RawAudClaim.([]interface{}) {
				claims.Audiences = append(claims.Audiences, a.(string))
			}
		}
	}
	return claims.Email, nil
}
