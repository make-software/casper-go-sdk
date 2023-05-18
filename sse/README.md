## SSE client

SSE package provide basic functionality to work with Casper events that streamed by SSE server.
It connects to the server and collect events to go channel, from other side consumers obtain this stream and delegate the process to specific handlers.

The example of simple usage:
```
func main() {
	client := sse.NewClient("http://52.3.38.81:9999/events/main")
	defer client.Stop()
	client.RegisterHandler(
		sse.DeployProcessedEventType,
		func(ctx context.Context, rawEvent sse.RawEvent) error {
			log.Printf("eventID: %d, raw data: %s", rawEvent.EventID, rawEvent.Data)
			deploy, err := rawEvent.ParseAsDeployProcessedEvent()
			if err != nil {
				return err
			}
			log.Printf("Deploy hash: %s", deploy.DeployProcessed.DeployHash)
			return nil
		})
	lastEventID := 1234
	client.Start(context.TODO(), lastEventID)
}
```

For more examples, please check [example_test.go](../tests/sse/example_test.go)

The SSE client is flexible and configurable, for the advanced usage check the [advanced doc](README_ADVANCED.md)