// Code generated by hack/generators/testcases-registration/main.go. DO NOT EDIT.
package integration

func init() {
	addTestsToTestSuite(
		TestAIGatewayCreation,
		TestControlPlaneEssentials,
		TestControlPlaneUpdate,
		TestControlPlaneWhenNoDataPlane,
		TestDataPlaneBlueGreenHorizontalScaling,
		TestDataPlaneBlueGreenRollout,
		TestDataPlaneBlueGreen_ResourcesNotDeletedUntilOwnerIsRemoved,
		TestDataPlaneEssentials,
		TestDataPlaneHorizontalScaling,
		TestDataPlaneUpdate,
		TestDataPlaneValidation,
		TestDataPlaneVolumeMounts,
		TestGatewayClassCreation,
		TestGatewayClassUpdates,
		TestGatewayConfigurationEssentials,
		TestGatewayDataPlaneNetworkPolicy,
		TestGatewayEssentials,
		TestGatewayMultiple,
		TestGatewayProvisionDataPlaneFail,
		TestGatewayWithMultipleListeners,
		TestHTTPRoute,
		TestHTTPRouteWithTLS,
		TestIngressEssentials,
		TestManualGatewayUpgradesAndDowngrades,
		TestScalingDataPlaneThroughGatewayConfiguration,
	)
}
