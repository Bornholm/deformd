var foo = params.get("foo")

if (foo === null) {
    throw new Error("params.get('foo') should not return null")
}

if (foo.bar !== 1) {
    throw new Error("foo.bar should be 1")
}

var miss = params.get("miss")

if (miss != null) {
    throw new Error("params.get('miss') should return null")
}