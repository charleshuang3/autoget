Please imp "/api/indexers", the spec is:

```
  /api/indexers:
    get:
      summary: Retrieve available Prowlarr indexers
      responses:
        '200':
          description: A list of indexers
          content:
            application/json:
              schema:
                type: object
                properties:
                  indexers:
                    type: array
                    items:
                      type: object
                      properties:
                        id:
                          type: integer
                        name:
                          type: string
```

Add Prowlarr structure in `backend/lib/handlers/prowlarr/prowlarr.go` which holds Config from `backend/lib/config/config.go`, so following func can use APIURL and APIKey from Config.

Impl the http handler in `backend/lib/handlers/prowlarr/indexers.go`, also add a func to make restful call prowlarr get indexers api in this file.

add router to `backend/lib/handlres/router.go` use `github.com/go-chi/chi/v5`.

