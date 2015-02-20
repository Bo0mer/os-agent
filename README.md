os-agent
========

HTTP Server acting as an OS agent written in Go

## Running it

### Configuration
Create configuration directory, which should contain `config.yml` file. You can preview example config file in the config/config.yml.

When you're done with the file, you should provide the config to the application by exporting a env variable `OS_AGENT_CONFIG_DIR` containing the full path to the directory, where your config.yml file is stored.

### Actual run
As simple as that:

```bash
go run main.go
```

## Running the tests
You need to have gingko installed.

```bash
ginkgo -r
```

## API

Following is the API that is provided by the OS-Agent for executing commands.

The response codes that are returned by the OS-Agent are splitted into the following groups:

| Group | Description |
| ---- | ----------- |
| 2XX | The requested operation completed successfully. |
| 4XX | There was a problem with the request payload. |
| 5XX | Execution of the operation failed due to some unexpected reason. |

**Payload**

The API is based on JSON request and responses. If not stated otherwise, default content-type should be `application/json`.

### Create Job

`POST /jobs`

**Request**

| Name | Type | Description |
| ---- | ---- | ----------- |
| async | boolean | Indicates whether the execution should be sync or async. |
| command | struct | Properties of the command. |

Example Request:

```JSON
{
    "async": true,
    "command": {
        "name": "cat",
        "args": [
            "arg1",
            "arg2"
        ],
        "env": {
            "variable_name1": "value1",
            "variable_name2": "value2"
        },
        "use_isolated_env": false,
        "working_dir": "/home/agent/whoa",
        "input": "This is the input to the cat command."
    }
}
```

**Response**

| Name | Type | Description |
| ---- | ---- | ----------- |
| id | string | The id of the job. |
| status | string | The status of the job. Either `IN_PROCESS` or `COMPLETED` |
| result | struct | The result of the command execution. |

Example Response:

```JSON
{
    "id": "jobid",
    "status": "COMPLETED",
    "result": {
        "stdout": "blabla",
        "stderr": "",
        "exit_code": 0,
        "error": ""
    }
}
```

### Get Job by Id
`GET /jobs?id=<job_id>`

**Request**

No data should be provided here.

**Response**


**Response**

| Name | Type | Description |
| ---- | ---- | ----------- |
| id | string | The id of the job. |
| status | string | The status of the job. Either `IN_PROCESS` or `COMPLETED` |
| result | struct | The result of the command execution. |

Example Response:

```JSON
{
    "id": "jobid",
    "status": "COMPLETED",
    "result": {
        "stdout": "blabla",
        "stderr": "",
        "exit_code": 0,
        "error": ""
    }
}
```
