package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mrtrkmnhub/cservice-cli/grpcconn"
	pb "github.com/mrtrkmnhub/cservice-cli/proto"
	"gopkg.in/yaml.v2"
)

var (
	CHALLENGE_TAG = os.Getenv("CHALLENGE_TAG")
	CAT_TAG       = os.Getenv("CAT_TAG") // assuming that categories are already in db.
)

func main() {

	port, _ := strconv.ParseUint(grpcconn.PORT, 0, 64)

	conf := grpcconn.Config{
		Endpoint: grpcconn.ENDPOINT,
		Port:     port,
		AuthKey:  grpcconn.AUTH_KEY,
		SignKey:  grpcconn.SIGN_KEY,
	}

	c, err := grpcconn.NewExServiceConn(conf)
	if err != nil {
		fmt.Printf("[exercise-service]  unable to connect %v", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var body interface{}
	fileContent, err := ioutil.ReadFile("./challenge-config/challenge.yml")
	if err != nil {
		log.Printf("ChallengeYAML err   #%v ", err)
	}
	fmt.Printf("Input: %s\n", fileContent)
	if err := yaml.Unmarshal(fileContent, &body); err != nil {
		panic(err)
	}

	body = convert(body)

	if b, err := json.Marshal(body); err != nil {
		panic(err)
	} else {
		fmt.Printf("Output: %s\n", b)
		_, err = c.AddExercise(ctx, &pb.AddExerciseRequest{
			Tag:         CHALLENGE_TAG,
			CategoryTag: CAT_TAG,
			Content:     string(b),
		})
		if err != nil {
			panic(err)
		}
		// slack notification can be integrated !
	}
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
