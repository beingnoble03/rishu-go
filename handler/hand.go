package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"
	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func Hand(c *gin.Context) {
	// configure cognito srp
	csrp, _ := cognitosrp.NewCognitoSRP("3de46ac0-2b2f-4847-ae42-4e17f8594cb8", "Test123@", "us-west-2_1BvlXo2Uu", "2vf9omb6ur3dmqtvv2bttb0h1l", nil)

	// configure cognito identity provider
	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
	)
	svc := cip.NewFromConfig(cfg)

	// initiate auth
	resp, err := svc.InitiateAuth(context.Background(), &cip.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	})
	if err != nil {
		panic(err)
	}

	// respond to password verifier challenge
	if resp.ChallengeName == types.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

		resp, err := svc.RespondToAuthChallenge(context.Background(), &cip.RespondToAuthChallengeInput{
			ChallengeName:      types.ChallengeNameTypePasswordVerifier,
			ChallengeResponses: challengeResponses,
			ClientId:           aws.String(csrp.GetClientId()),
		})
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"authParam": csrp.GetAuthParams(),
		})
		// print the tokens
		fmt.Printf("Access Token: %s\n", *resp.AuthenticationResult.AccessToken)
		fmt.Printf("ID Token: %s\n", *resp.AuthenticationResult.IdToken)
		fmt.Printf("Refresh Token: %s\n", *resp.AuthenticationResult.RefreshToken)
		fmt.Println("AuthParam")
		fmt.Println(csrp.GetAuthParams())
	}
}
