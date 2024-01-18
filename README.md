# applife
Run and gracefully terminate processes

```
// Create app instance. Supply optional logger struct for feedback on what is happening.
app := applife.NewApp(ctx, new(logger))

// Add processes to app which will execute once Run is called.
app.AddProcess("process_1", func(ctx context.Context) {
  // Process 1 code here
}
app.AddProcess("process_2", func(ctx context.Context) {
  // Process 2 code here
}

// Run the processes. This is a blocking call.
app.Run()

// Context will be cancelled once terminate or interrupt signal is received
```
