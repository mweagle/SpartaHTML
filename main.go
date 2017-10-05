package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
	spartaCF "github.com/mweagle/Sparta/aws/cloudformation"
)

////////////////////////////////////////////////////////////////////////////////
// Hello world event handler
//
func helloWorld(w http.ResponseWriter, r *http.Request) {
	logger, _ := r.Context().Value(sparta.ContextKeyLogger).(*logrus.Logger)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var jsonMessage json.RawMessage
	err := decoder.Decode(&jsonMessage)
	if err != nil {
		logger.Error("Failed to decode request: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("Hello World: ", string(jsonMessage))
	w.Write(jsonMessage)
}

func spartaLambdaFunctions(api *sparta.API) []*sparta.LambdaAWSInfo {
	var lambdaFunctions []*sparta.LambdaAWSInfo
	lambdaFn := sparta.HandleAWSLambda(sparta.LambdaName(helloWorld),
		http.HandlerFunc(helloWorld),
		sparta.IAMRoleDefinition{})

	if nil != api {
		apiGatewayResource, _ := api.NewResource("/hello", lambdaFn)
		_, err := apiGatewayResource.NewMethod("GET", http.StatusOK)
		if nil != err {
			panic("Failed to create /hello resource")
		}
	}
	return append(lambdaFunctions, lambdaFn)
}

////////////////////////////////////////////////////////////////////////////////
// Main
func main() {
	// Register the function with the API Gateway
	apiStage := sparta.NewStage("v1")
	apiGateway := sparta.NewAPIGateway("SpartaHTML", apiStage)
	// Enable CORS s.t. the S3 site can access the resources
	apiGateway.CORSEnabled = true

	// Provision a new S3 bucket with the resources in the supplied subdirectory
	s3Site, _ := sparta.NewS3Site("./resources")

	// Deploy it
	stackName := spartaCF.UserScopedStackName("SpartaHTML")
	sparta.Main(stackName,
		fmt.Sprintf("Sparta app that provisions a CORS-enabled API Gateway together with an S3 site"),
		spartaLambdaFunctions(apiGateway),
		apiGateway,
		s3Site)
}
