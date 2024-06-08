framework/plugins 是以作为"插件"形式存在的, 包含预选、优选.

与之相对的, 是调度器初始化时加载的内置算法.

pkg/scheduler/algorithmprovider/defaults/register_predicates.go -> init()
pkg/scheduler/algorithmprovider/defaults/register_priorities.go -> init()
