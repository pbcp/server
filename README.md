> **Warning:** Currently under volatile development; not ready for general use!

## pbcp

The server.

### Dependencies

- [boltdb/bolt](https://github.com/boltdb/bolt)
- [aws/aws-sdk-go](https://github.com/aws/aws-sdk-go)
- [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)

### Configuration

- The database is currently written to a `db/` subdirectory; make sure this exists and is writable by the running user
- Clips are stored in S3
    - See [the `aws-sdk-go` wiki](https://github.com/aws/aws-sdk-go/wiki/configuring-sdk#specifying-credentials) for instructions on setting up credentials
    - Currently, a bucket `pbcp` must exist in the `us-west-2` region
    - The bucket policy should be set to allow anonymous reads (for file retrieval)

### API

- `GET /board/:id/:num`: Get clip at index (starts from 0)
- `POST /board/:id`: Upload file to clipboard (request body should be file body)
- `GET /register`: Generate a unique user ID to use
