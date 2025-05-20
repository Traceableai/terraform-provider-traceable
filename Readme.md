# Terraform Provider Testing Guide

## How to Write and Run Unit Tests

### 1. Create a Test File
- Place the test file under the `resources` folder.

### 2. Set Environment Variables
Before running tests, export the required environment variables:

```bash
export API_TOKEN=<your_token>
export PLATFORM_URL=<your_url>
```

### 3. Write a Basic Test
A basic test should create, import, update, and delete the resource. Here's an example:

```go
resource.Test(t, resource.TestCase{
    PreCheck:                 func() { acctest.TestAccPreCheck(t) },
    ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
    Steps: []resource.TestStep{
        {
            Config: testAccRateLimitiningResourceConfigDefault("rate_limit_T1", "ALERT"),
            Check: resource.ComposeAggregateTestCheckFunc(
                resource.TestCheckResourceAttr("traceable_rate_limiting.test", "name", "rate_limit_T1"),
                resource.TestCheckResourceAttr("traceable_rate_limiting.test", "action.action_type", "ALERT"),
                resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NU"),
                resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NF"),
            ),
        },
        {
            ResourceName:      "traceable_rate_limiting.test",
            ImportState:       true,
            ImportStateId:     "rate_limit_T1",
            ImportStateVerify: true,
        },
        {
            Config: testAccRateLimitiningResourceConfigDefault("rate_limit_T1", "BLOCK"),
            Check: resource.ComposeAggregateTestCheckFunc(
                resource.TestCheckResourceAttr("traceable_rate_limiting.test", "name", "rate_limit_T1"),
                resource.TestCheckResourceAttr("traceable_rate_limiting.test", "action.action_type", "BLOCK"),
                resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NU"),
                resource.TestCheckTypeSetElemAttr("traceable_rate_limiting.test", "sources.regions.regions_ids.*", "NF"),
            ),
        },
        // Delete testing automatically occurs in TestCase
    },
})
```

### 4. Run Tests
Run this command to test your changes. Replace the file name:

```bash
TF_ACC=5 go test -v internal/resources/resource_rate_limiting_test.go
```

Or to test all files:

```bash
TF_ACC=5 go test -v ./...
```

### 5. Test Steps
In every test step, there is an implicit check that Terraform plan should be empty after each apply. Otherwise, the test will fail.
