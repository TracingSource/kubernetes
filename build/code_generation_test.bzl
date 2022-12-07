load(":code_generation.bzl", "bazel_go_library", "go_pkg")
load("@bazel_skylib//lib:unittest.bzl", "asserts", "unittest")

def _bazel_go_library_test_impl(ctx):
    env = unittest.begin(ctx)
    test_cases = [
        ("pkg/kubectl/util", "//pkg/kubectl/util:go_default_library"),
        ("vendor/some/third/party", "//vendor/some/third/party:go_default_library"),
        ("staging/src/k8s.io/apimachinery/api", "//staging/src/k8s.io/apimachinery/api:go_default_library"),
    ]
    for input, expected in test_cases:
        asserts.equals(env, expected, bazel_go_library(input))
    unittest.end(env)

bazel_go_library_test = unittest.make(_bazel_go_library_test_impl)

def _go_pkg_test_impl(ctx):
    env = unittest.begin(ctx)
    test_cases = [
        ("pkg/kubectl/util", "k8s.io/kubernetes/pkg/kubectl/util"),
        ("vendor/some/third/party", "k8s.io/kubernetes/vendor/some/third/party"),
        ("staging/src/k8s.io/apimachinery/api", "k8s.io/kubernetes/vendor/k8s.io/apimachinery/api"),
    ]
    for input, expected in test_cases:
        asserts.equals(env, expected, go_pkg(input))
    unittest.end(env)

go_pkg_test = unittest.make(_go_pkg_test_impl)

def code_generation_test_suite(name):
    unittest.suite(
        name,
        bazel_go_library_test,
        go_pkg_test,
    )
