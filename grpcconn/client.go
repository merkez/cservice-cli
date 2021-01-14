package grpcconn

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/mrtrkmnhub/cservice-cli/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	UnreachableErr = errors.New("seems to be unreachable")

	NoTokenErrMsg     = "token contains an invalid number of segments"
	UnauthorizedErr   = errors.New("Unauthorized attempt to use exercise service ")
	UnauthorizeErrMsg = "unauthorized"
	AUTH_KEY          = os.Getenv("AUTH_KEY")
	SIGN_KEY          = os.Getenv("SIGN_KEY")
	ENDPOINT          = os.Getenv("ENDPOINT")
	PORT              = os.Getenv("PORT")
)

type Config struct {
	Endpoint string
	Port     uint64
	AuthKey  string
	SignKey  string
}

type Creds struct {
	Token    string
	Insecure bool
}

func (c Creds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": string(c.Token),
	}, nil
}

func (c Creds) RequireTransportSecurity() bool {
	return !c.Insecure
}

// only secure communication on CI
func NewExServiceConn(config Config) (pb.ExerciseStoreClient, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"au": config.AuthKey,
	})
	tokenString, err := token.SignedString([]byte(config.SignKey))
	if err != nil {
		return nil, TranslateRPCErr(err)
	}

	authCreds := Creds{Token: tokenString}

	pool, _ := x509.SystemCertPool()
	creds := credentials.NewClientTLSFromCert(pool, "")

	creds = credentials.NewTLS(&tls.Config{
		RootCAs: pool,
	})

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(authCreds),
	}

	conn, err := grpc.Dial(config.Endpoint+":"+strconv.FormatUint(config.Port, 10), dialOpts...)
	if err != nil {
		log.Error().Msgf("Error on dialing vpn service: %v", err)
		return nil, TranslateRPCErr(err)
	}
	c := pb.NewExerciseStoreClient(conn)
	return c, nil

}

func TranslateRPCErr(err error) error {
	st, ok := status.FromError(err)
	if ok {
		msg := st.Message()
		switch {
		case UnauthorizeErrMsg == msg:
			return UnauthorizedErr

		case NoTokenErrMsg == msg:
			return UnauthorizedErr

		case strings.Contains(msg, "TransientFailure"):
			return UnreachableErr
		}

		return err
	}

	return err
}
