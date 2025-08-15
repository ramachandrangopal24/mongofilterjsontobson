The filter package provides utilities to parse JSON request bodies into MongoDB-compatible BSON filters, supporting both MongoDB Extended JSON (preferred) and plain JSON.
It’s designed for use in HTTP services that need to pass client-provided query filters directly to MongoDB.

Features

Extended JSON parsing
Fully supports MongoDB Extended JSON syntax for special BSON types ($oid, $date, $in, $gte, etc.).

Plain JSON fallback
Automatically converts _id fields that look like 24-character hex strings into primitive.ObjectID.

Optional query wrapper parsing
Supports a structured request with filter, sort, projection, limit, and skip fields for advanced querying.

Order-preserving BSON conversion
Outputs bson.D to preserve field order when required.
