# API Endpoints

- POST /log/error/list
  - Description: List error logs
  - JSON Request:
    - "app": (string)
    - "page": (Paging)
      - "limit": (int) page limit
      - "page": (int) page number, 1-based
      - "total": (int) total count
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
  - cURL:
    ```sh
    curl -X POST 'http://localhost:8087/log/error/list' \
      -H 'Content-Type: application/json' \
      -d '{"app":"","page":{"limit":0,"page":0,"total":0}}'
    ```

  - JSON Request Object In TypeScript:
    ```ts
    export interface ListErrorLogReq {
      app?: string
      page?: Paging
    }
    export interface Paging {
      limit?: number                 // page limit
      page?: number                  // page number, 1-based
      total?: number                 // total count
    }
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface Resp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    let req: ListErrorLogReq | null = null;
    this.http.post<Resp>(`/log/error/list`, req)
      .subscribe({
        next: (resp: Resp) => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /auth/resource
  - Description: Expose resource and endpoint information to other backend service for authorization.
  - Expected Access Scope: PROTECTED
  - JSON Response:
    - "errorCode": (string) error code
    - "msg": (string) message
    - "error": (bool) whether the request was successful
    - "data": (ResourceInfoRes) response data
      - "resources": ([]auth.Resource)
        - "name": (string) resource name
        - "code": (string) resource code, unique identifier
      - "paths": ([]auth.Endpoint)
        - "type": (string) access scope type: PROTECTED/PUBLIC
        - "url": (string) endpoint url
        - "group": (string) app name
        - "desc": (string) description of the endpoint
        - "resCode": (string) resource code
        - "method": (string) http method
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/auth/resource'
    ```

  - JSON Response Object In TypeScript:
    ```ts
    export interface GnResp {
      errorCode?: string             // error code
      msg?: string                   // message
      error?: boolean                // whether the request was successful
      data?: ResourceInfoRes
    }
    export interface ResourceInfoRes {
      resources?: Resource[]
      paths?: Endpoint[]
    }
    export interface Resource {
      name?: string                  // resource name
      code?: string                  // resource code, unique identifier
    }
    export interface Endpoint {
      type?: string                  // access scope type: PROTECTED/PUBLIC
      url?: string                   // endpoint url
      group?: string                 // app name
      desc?: string                  // description of the endpoint
      resCode?: string               // resource code
      method?: string                // http method
    }
    ```

  - Angular HttpClient Demo:
    ```ts
    this.http.get<GnResp>(`/auth/resource`)
      .subscribe({
        next: (resp: GnResp) => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /metrics
  - Description: Collect prometheus metrics information
  - Header Parameter:
    - "Authorization": Basic authorization if enabled
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/metrics' \
      -H 'Authorization: '
    ```

  - Angular HttpClient Demo:
    ```ts
    let authorization: any | null = null;
    this.http.get<any>(`/metrics`,
      {
        headers: {
          "Authorization": authorization
        }
      })
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /debug/pprof
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/debug/pprof'
    ```

  - Angular HttpClient Demo:
    ```ts
    this.http.get<any>(`/debug/pprof`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /debug/pprof/:name
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/debug/pprof/:name'
    ```

  - Angular HttpClient Demo:
    ```ts
    this.http.get<any>(`/debug/pprof/:name`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /debug/pprof/cmdline
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/debug/pprof/cmdline'
    ```

  - Angular HttpClient Demo:
    ```ts
    this.http.get<any>(`/debug/pprof/cmdline`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /debug/pprof/profile
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/debug/pprof/profile'
    ```

  - Angular HttpClient Demo:
    ```ts
    this.http.get<any>(`/debug/pprof/profile`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /debug/pprof/symbol
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/debug/pprof/symbol'
    ```

  - Angular HttpClient Demo:
    ```ts
    this.http.get<any>(`/debug/pprof/symbol`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /debug/pprof/trace
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/debug/pprof/trace'
    ```

  - Angular HttpClient Demo:
    ```ts
    this.http.get<any>(`/debug/pprof/trace`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```

- GET /doc/api
  - Description: Serve the generated API documentation webpage
  - Expected Access Scope: PUBLIC
  - cURL:
    ```sh
    curl -X GET 'http://localhost:8087/doc/api'
    ```

  - Angular HttpClient Demo:
    ```ts
    this.http.get<any>(`/doc/api`)
      .subscribe({
        next: () => {
        },
        error: (err) => {
          console.log(err)
        }
      });
    ```
