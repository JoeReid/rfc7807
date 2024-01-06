RFC7807
====

```golang
// MyHandler is an http handler that uses the rfc7807 lib to return nicely formatted errors
func MyHandler(w http.ResponseWriter, r *http.Request) {
    data, err := GetDataFromDB()
    if err != nil {
        rfc7807.ErrorResponse(err, w, r)
        return
    }

    // Continue to process data and generate a response
}

// GetDataFromDB returns an error as an example. Real code would have many execution paths
func GetDataFromDB() (MyDataType, error) {
    return nil, &rfc7807.Error{
        Cause: err,
        Status: http.StatusInternalServerError,
        Type: "https://example.net/errors/internal-server-error",
        Detail: "While handling the request, the database connection failed or was never established",
    }
}
```