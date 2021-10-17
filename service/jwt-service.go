package service

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/br3w0r/gamelist-backend/entity"
	"github.com/br3w0r/gamelist-backend/repository"
	utilErrs "github.com/br3w0r/gamelist-backend/util/errors"
	"github.com/golang-jwt/jwt"
)

type JWTService interface {
	GenerateTokens(user string) (*entity.TokenPair, error)
	Authenticate(tokenString string) (string, error)
	RefreshTokens(refreshToken string) (*entity.TokenPair, error)
	RevokeRefreshToken(refreshToken string) error
	DeleteAllUserRefreshTokens(nickname string) error
}

type jwtService struct {
	repo              repository.GamelistRepository
	secretKeys        []string
	refreshSecretKeys []string
}

func NewJWTService(repo repository.GamelistRepository) JWTService {
	return &jwtService{
		repo: repo,
		secretKeys: []string{
			"-U=cp#@Frd/kVYM(5nVNHtW27;;ptGkpI*TFPWshC1h@BFBS!f;?nE+w;uAw+:c6tV4P7g",
			"hAHsjr#IT5>f3H+MwT?%oMPBV&jVq,gk,Q/5O!b1db'o&!DxpAIs%dzt<j-c5#qP:e?-xz",
			"19zJELBzJqC7HFLWqJQJ/Vez6/z)I49'hzyzCzV&nHZ6+E<zHVv'P)28+6!?=8*09vr328",
			"xPz:btyjLvB)t%,<R%Vk@-NUwzq%V=K5<WGgjZ3%+zRpf4V)m-Y:HuuH1SP-=J.dl'4+wE",
			"=Ei:/D4t)Rha*FJ,nrd7@L*wiu?b,-lMt@3AYaq06EMhvdd@;vYcF!RjaOrdkq=z8*45P+",
			">u8b1DS3cg8R7L:R0!9bx.uN?q.gE#.WV7+H*O;OHM.BGR9TZ%@rlRR*th&/qXL'LGmmkT",
			"31=Dsa5U@D(5C+S#pN<8M8C9OJzG'=FR(Pigyvda>rml)7C$Aon(,gjXz09yfgNs1:<(Q8",
			"n1-NTiT#peNRDfi-6Cj26P9>=@6!hXJW<FtPpQH$F2rBYZ00>*nAhB;%5wV??aF8fzc?&/",
			"KPK?y5EFC($aCLWIYJFzZQtd&6)DcM#Ft+zwuXt4xcA>lGV#bvvaX?',2o3A!<Yv7Ay,J$",
			"Kjv(TUAXXIEntY1Fz)So<;jX+.$Dsu%K<8@6:04<1=X3:hb'U=iXy=2wEFobCJS4r*Nb5>",
		},
		refreshSecretKeys: []string{
			"1d2Y:5fbjoE.O@AV,SLg6!6yE(@UyTylqy&ag/T51bvkc#o3jW5JsP:6P@naaLt-Dd6iPz",
			"j@@mUI+5C),mnK:n@JTlVfT$IXOb/m6t1J@Dg3k;L@b*(zw#aD)spFIr5iSy15Gj(&&.dP",
			"&$pg?wMrWdh%dY6VlPr+Pd(<0Pk:tCNp8qr0s0KGQuGZ.;r)Cx@PB,?1qcMj-(aE6FCC$)",
			"mrd=Q4D:m*LWcssCwCwbSBa?F-nNpr%Yols'mk<R!eyx*2qM#$097syZUiK3H'GCCIj7EZ",
			"uPetcuyAu@x'XyFHY5LdXJeWo+'4E4Y<ZW0EYJl27MSE1,?cF71>xK3Ao8w:dz:JZcAN-A",
			"PsrND<?jILG@q=c87lv9Xc(A+-g!4cK7B/N5a;s1cH*wS'IYbjtgu?$SoaO1aDL&!1cVU=",
			"&NNA!im9h(G(DI$ruQF+t;d/Ohtiif/*yT&QMc9c#UMXQZSYaFm!v-2$<O%Es:Gr5E@q73",
			"hm@+o6gb7EZaC)p:,S=YtrW7re*=Sk6CWXe0fe5!&W5O;91cpAt3w+f7h,+CqgT0F*<S?6",
			"*aSa*ybrI18&JDkKZMz2dhESIVOA$$CRcwg%V%iA.D+)Ol.@aLb(6dD<41fazyyk&Gl-N4",
			"e73GiGg?:qqg;K'*8YaK7u,YC:yz9v'gF)y'NGqfY+a*vL4V%V/qv>y&?z11ODf'd+y6'0",
		},
	}
}

func (s *jwtService) GenerateTokens(user string) (*entity.TokenPair, error) {
	iat := time.Now().Unix()
	exp := iat + 3900

	kid := rand.Intn(10)
	key := []byte(s.secretKeys[kid] + "." + strconv.FormatInt(iat, 10))
	claims := jwt.MapClaims{
		"sub": user,
		"iat": iat,
		"exp": exp,
		"kid": kid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return nil, utilErrs.New(utilErrs.Internal, err, "failed to generate token string")
	}

	refreshKid := rand.Intn(10)
	refreshKey := []byte(s.refreshSecretKeys[refreshKid] + "." + strconv.FormatInt(iat, 10))
	refreshClaims := jwt.MapClaims{
		"sub": user,
		"iat": iat,
		"kid": refreshKid,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return nil, utilErrs.New(utilErrs.Internal, err, "failed to generate refresh token string")
	}

	err = s.repo.SaveRefreshToken(user, refreshTokenString)
	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *jwtService) Authenticate(tokenString string) (string, error) {
	return s.validateToken(tokenString, false)
}

func (s *jwtService) RefreshTokens(refreshToken string) (*entity.TokenPair, error) {
	user, err := s.validateToken(refreshToken, true)
	if err != nil {
		return nil, err
	}

	err = s.repo.FindRefreshToken(user, refreshToken)
	if err != nil {
		return nil, err
	}

	tokens, err := s.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	err = s.repo.DeleteRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *jwtService) RevokeRefreshToken(refreshToken string) error {
	return s.repo.DeleteRefreshToken(refreshToken)
}

func (s *jwtService) DeleteAllUserRefreshTokens(nickname string) error {
	return s.repo.DeleteAllUserRefreshTokens(nickname)
}

func (s *jwtService) validateToken(tokenString string, isRefresh bool) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utilErrs.Newf(utilErrs.BadInput, nil, "invalid signing method: %v", token.Header["alg"])
		}
		if err := token.Claims.Valid(); err != nil {
			return nil, utilErrs.New(utilErrs.BadInput, err, "failed to validate claims")
		}

		claims := token.Claims.(jwt.MapClaims)

		var secret string
		if isRefresh {
			secret = s.refreshSecretKeys[int(claims["kid"].(float64))]
		} else {
			secret = s.secretKeys[int(claims["kid"].(float64))]
		}

		key := []byte(secret + "." + strconv.FormatInt(int64(claims["iat"].(float64)), 10))
		return key, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}

	if !token.Valid {
		return "", utilErrs.New(utilErrs.Unauthorized, nil, "token validation failed")
	}

	return "", err
}
