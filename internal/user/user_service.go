package user

import (
	"context"
	// "log"
	"os"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type service struct {
    Repository
    timeout time.Duration
}

func NewService(repository Repository) Service {
    return &service{
        repository,
        time.Duration(2) * time.Second,
    }
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
    ctx, cancel := context.WithTimeout(c, s.timeout)
    defer cancel()

    hashedPassword, err := util.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    // Log hashed password yang dihasilkan
    // log.Printf("Generated hashed password: %s", hashedPassword)

    u := &User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
    }

    r, err := s.Repository.CreateUser(ctx, u)
    if err != nil {
        return nil, err
    }

    res := &CreateUserRes{
        ID:       strconv.Itoa(int(r.ID)),
        Username: r.Username,
        Email:    r.Email,
    }

    return res, nil
}

type JWTClaims struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
    ctx, cancel := context.WithTimeout(c, s.timeout)
    defer cancel()

    u, err := s.Repository.GetUserByEmail(ctx, req.Email)
    if err != nil {
        return &LoginUserRes{}, err
    }

    // Log stored hashed password
    // log.Printf("Stored hashed password: %s", u.Password)
    // Log input password
    // log.Printf("Input password: %s", req.Password)

    err = util.CheckPassword(req.Password, u.Password)
    if err != nil {
        return &LoginUserRes{}, err
    }

    secretKey := os.Getenv("SECRET_KEY")
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
        ID:       strconv.Itoa(int(u.ID)),
        Username: u.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    strconv.Itoa(int(u.ID)),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    })

    ss, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return &LoginUserRes{}, err
    }

    return &LoginUserRes{accessToken: ss, Username: u.Username, ID: strconv.Itoa(int(u.ID))}, nil
}
