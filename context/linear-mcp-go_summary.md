This project is a **Model Context Protocol (MCP) server for Linear**, written in **Go**. Its primary goal is to expose the Linear API as a set of tools that AI assistants can use to interact with Linear projects. This allows users to perform actions like creating, updating, and searching for issues directly from their AI-powered development environment.

### Technology and Approach

The project is built using a modern Go stack and follows best practices for creating robust and maintainable software.

*   **Go:** The core language used for the server, providing high performance and strong typing.
*   **Cobra:** A popular library for building command-line interfaces in Go. It's used to create the `serve` and `setup` commands.
*   **MCP (Model Context Protocol):** The server implements the MCP, which defines a standardized way for AI assistants to discover and use external tools.
*   **Linear API:** The server acts as a bridge to the Linear API, handling authentication and rate limiting.
*   **GitHub Actions:** The project uses GitHub Actions for continuous integration and automated releases.

### Test Approach

The testing strategy is a key aspect of this project, designed to ensure the reliability of the Linear integration without depending on a live API for every test run.

*   **`go-vcr` for API Mocking:** The project uses the `go-vcr` library to record and replay HTTP interactions with the Linear API. This approach has several advantages:
    *   **Offline Testing:** Tests can be run without an internet connection, making them faster and more reliable.
    *   **Deterministic Tests:** The recorded responses (cassettes) ensure that tests always run against the same data, eliminating flakiness due to API changes or network issues.
    *   **Isolation:** The tests are isolated from the actual Linear project, so they don't create or modify real data during test runs.
*   **Cassettes and Golden Files:**
    *   **Cassettes:** The recorded HTTP interactions are stored in YAML files (cassettes) in the `testdata/fixtures` directory. Each cassette corresponds to a specific test case and contains the request and response data.
    *   **Golden Files:** The expected output of each test is stored in `.golden` files in the `testdata/golden` directory. During a test run, the actual output is compared to the golden file. If they don't match, the test fails.
*   **Test Workflows:**
    *   **Running Tests:** `go test -v ./...` runs the tests against the recorded cassettes.
    *   **Re-recording Cassettes:** `go test -v -record=true ./...` re-records the cassettes by making live requests to the Linear API. This is useful when the API changes or new tests are added.
    *   **Updating Golden Files:** `go test -v -golden=true ./...` updates the golden files with the latest output from the tests.

This comprehensive testing approach, combining API mocking with golden file testing, ensures that the Linear MCP server is both correct and resilient. It allows developers to confidently make changes to the codebase, knowing that the tests will catch any regressions in the integration with the Linear API.
