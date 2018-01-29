package main

import (
	"context"
	"fmt"
	"net/http"

	sparta "github.com/mweagle/Sparta"
	spartaCF "github.com/mweagle/Sparta/aws/cloudformation"
	spartaAWSEvents "github.com/mweagle/Sparta/aws/events"
	gocf "github.com/mweagle/go-cloudformation"
	"github.com/sirupsen/logrus"
)

type helloWorldResponse struct {
	Message string
	Request spartaAWSEvents.APIGatewayRequest
}

////////////////////////////////////////////////////////////////////////////////
// Hello world event handler
func helloWorld(ctx context.Context,
	gatewayEvent spartaAWSEvents.APIGatewayRequest) (interface{}, error) {
	/*
		 To return an error back to the client using a standard HTTP status code:

			errorResponse := spartaAPIG.NewErrorResponse(http.StatusInternalError,
			"Something failed inside here")
			return errorResponse, nil

			You can also create custom error response types, so long as they
			include `"code":HTTP_STATUS_CODE` somewhere in the response body.
			This reserved expression is what Sparta uses as a RegExp to determine
			the Integration Mapping value
	*/

	logger, loggerOk := ctx.Value(sparta.ContextKeyLogger).(*logrus.Logger)
	if loggerOk {
		logger.Info("Hello world structured log message")
	}

	// Return a message, together with the incoming input...
	return &helloWorldResponse{
		Message: fmt.Sprintf("Hello world 🌏"),
		Request: gatewayEvent,
	}, nil
}

func spartaHTMLLambdaFunctions(api *sparta.API) []*sparta.LambdaAWSInfo {
	var lambdaFunctions []*sparta.LambdaAWSInfo
	lambdaFn := sparta.HandleAWSLambda(sparta.LambdaName(helloWorld),
		helloWorld,
		sparta.IAMRoleDefinition{})

	if nil != api {
		apiGatewayResource, _ := api.NewResource("/hello", lambdaFn)

		// We only return http.StatusOK
		apiMethod, apiMethodErr := apiGatewayResource.NewMethod("GET",
			http.StatusOK,
			http.StatusInternalServerError)
		if nil != apiMethodErr {
			panic("Failed to create /hello resource: " + apiMethodErr.Error())
		}
		// The lambda resource only supports application/json Unmarshallable
		// requests.
		apiMethod.SupportedRequestContentTypes = []string{"application/json"}
	}
	return append(lambdaFunctions, lambdaFn)
}

////////////////////////////////////////////////////////////////////////////////
// Main
func main() {
	// Provision an S3 site
	s3Site, s3SiteErr := sparta.NewS3Site("./resources")
	if s3SiteErr != nil {
		panic("Failed to create S3 Site")
	}

	// Register the function with the API Gateway
	apiStage := sparta.NewStage("v1")
	apiGateway := sparta.NewAPIGateway("SpartaHTML", apiStage)
	// Enable CORS s.t. the S3 site can access the resources
	apiGateway.CORSOptions = &sparta.CORSOptions{
		Headers: map[string]interface{}{
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Origin":  gocf.GetAtt(s3Site.CloudFormationS3ResourceName(), "WebsiteURL"),
		},
	}

	// Deploy it
	stackName := spartaCF.UserScopedStackName("SpartaHTML")
	sparta.Main(stackName,
		fmt.Sprintf("SpartaHTML provisions a static S3 hosted website with an API Gateway resource backed by a custom Lambda function"),
		spartaHTMLLambdaFunctions(apiGateway),
		apiGateway,
		s3Site)
}
