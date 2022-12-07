# Genrule wrapper around the go-bindata utility.
# IMPORTANT: Any changes to this rule may also require changes to hack/generate-bindata.sh.
def go_bindata(
    name, srcs, outs,
    compress=True,
    include_metadata=True,
    pkg="generated",
    ignores=["\.jpg", "\.png", "\.md", "BUILD(\.bazel)?"],
    **kw):

  args = []
  for ignore in ignores:
    args.extend(["-ignore", "'%s'" % ignore])
  if not include_metadata:
    args.append("-nometadata")
  if not compress:
    args.append("-nocompress")

  native.genrule(
    name = name,
    srcs = srcs,
    outs = outs,
    cmd = """
    $(location //vendor/github.com/go-bindata/go-bindata/go-bindata:go-bindata) \
      -o "$@" -pkg %s -prefix $$(pwd) %s $(SRCS)
    """ % (pkg, " ".join(args)),
    tools = [
      "//vendor/github.com/go-bindata/go-bindata/go-bindata",
    ],
    **kw
  )
