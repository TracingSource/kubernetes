// Package ipam provides different allocators for assigning IP ranges to nodes.
// We currently support several kinds of IPAM allocators (these are denoted by
// the CIDRAllocatorType):
// - RangeAllocator is an allocator that assigns PodCIDRs to nodes and works
//   in conjunction with the RouteController to configure the network to get
//   connectivity.
// - CloudAllocator is an allocator that synchronizes PodCIDRs from IP
//   ranges assignments from the underlying cloud platform.
// - (Alpha only) IPAMFromCluster is an allocator that has the similar
//   functionality as the RangeAllocator but also synchronizes cluster-managed
//   ranges into the cloud platform.
// - (Alpha only) IPAMFromCloud is the same as CloudAllocator (synchronizes
//   from cloud into the cluster.)
package ipam
