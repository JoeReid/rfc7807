RFC7807
====

```golang
func MyHandler(w http.ResponseWriter, r *http.Request) {
    if err := someLogicThatErrors(); err != nil {
        rfc7807.ErrorResponse(&rfc7807.Error{
            Cause:  err,
			Status: http.StatusBadRequest,
			Type:   "https://example.net/errors/bad-request",
			Detail: "Some custom details about the error",
        }, w, r)

        return
    }

    w.WriteHeader(http.StatusOK)
}
```
