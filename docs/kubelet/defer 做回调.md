# defer 做回调

```golang
// caller: pkg/volume/util/operationexecutor/operation_executor.go -> operationExecutor.MountVolume()
func (grm *nestedPendingOperations) Run(
	volumeName v1.UniqueVolumeName,
	podName types.UniquePodName,
	generatedOperations types.GeneratedOperations,
) error {
    // ...
	go func() (eventErr, detailedErr error) {
		// Handle unhandled panics (very unlikely)
		defer k8sRuntime.HandleCrash()
		// Handle completion of and error, if any, from operationFunc()
		// 当此 go 协程执行完成后, 进行一个回调.
		// 这种用 defer 做回调的方式, 真是巧妙.
		defer grm.operationComplete(volumeName, podName, &detailedErr)
		return generatedOperations.Run()
	}()

	return nil
}
```
