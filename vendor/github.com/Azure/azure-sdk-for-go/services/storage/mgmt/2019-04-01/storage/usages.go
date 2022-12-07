package storage

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// UsagesClient is the the Azure Storage Management API.
type UsagesClient struct {
	BaseClient
}

// NewUsagesClient creates an instance of the UsagesClient client.
func NewUsagesClient(subscriptionID string) UsagesClient {
	return NewUsagesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewUsagesClientWithBaseURI creates an instance of the UsagesClient client.
func NewUsagesClientWithBaseURI(baseURI string, subscriptionID string) UsagesClient {
	return UsagesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// ListByLocation gets the current usage count and the limit for the resources of the location under the subscription.
// Parameters:
// location - the location of the Azure Storage resource.
func (client UsagesClient) ListByLocation(ctx context.Context, location string) (result UsageListResult, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/UsagesClient.ListByLocation")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("storage.UsagesClient", "ListByLocation", err.Error())
	}

	req, err := client.ListByLocationPreparer(ctx, location)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.UsagesClient", "ListByLocation", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByLocationSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "storage.UsagesClient", "ListByLocation", resp, "Failure sending request")
		return
	}

	result, err = client.ListByLocationResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storage.UsagesClient", "ListByLocation", resp, "Failure responding to request")
	}

	return
}

// ListByLocationPreparer prepares the ListByLocation request.
func (client UsagesClient) ListByLocationPreparer(ctx context.Context, location string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"location":       autorest.Encode("path", location),
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-04-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Storage/locations/{location}/usages", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByLocationSender sends the ListByLocation request. The method will close the
// http.Response Body if it receives an error.
func (client UsagesClient) ListByLocationSender(req *http.Request) (*http.Response, error) {
	sd := autorest.GetSendDecorators(req.Context(), azure.DoRetryWithRegistration(client.Client))
	return autorest.SendWithSender(client, req, sd...)
}

// ListByLocationResponder handles the response to the ListByLocation request. The method always
// closes the http.Response Body.
func (client UsagesClient) ListByLocationResponder(resp *http.Response) (result UsageListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
