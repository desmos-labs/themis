import cffi

try:
    # These incantations are necessary to make ripemd160 work on OpenSSL 3
    ffi = cffi.FFI()
    ffi.cdef("void *OSSL_PROVIDER_load(void *libctx, const char *name);")
    lib = ffi.dlopen("crypto")
    assert (
            lib.OSSL_PROVIDER_load(
                ffi.NULL, ffi.new("char[]", "default".encode("ascii"))
            )
            is not ffi.NULL
    )
    assert (
            lib.OSSL_PROVIDER_load(
                ffi.NULL, ffi.new("char[]", "legacy".encode("ascii"))
            )
            is not ffi.NULL
    )
except Exception:
    # Ignore all failures, the code above is fragile and will fail in
    # situations when it is not needed (e.g. you don't have OpenSSL 3).
    pass
