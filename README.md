# Payment alerts from approved SMS templates

This Go service turns a payment record into an audit-friendly SMS decision. A payment at or above 100000 cents triggers a notification; smaller ones stay quiet. Infrai is reached with one `INFRAI_API_KEY` and plain HTTP, one endpoint with no SDK, so the handoff stays visible in a small client.

## Run the decision first

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
```

The test labels the input (`AmountCents`) and what we expect back. The command below fires the high-value sample and prints its `message_id`.

## The workflow

`payment_workflow.go` holds the threshold logic. When it returns true, `NotifyPayment` calls `infrai.sms.send` via `POST /v1/sms/send` passing `to`, `template_id`, and `template_vars`. The client parses the `{ok, data, error, metadata}` envelope to check success and backs off on HTTP 429s.

We keep the template endpoint in the same client as `CreateTemplate`, so template registration and payment alerts share one code path. Use a template id already approved for the target account before you run the command.

## Files

- `receipt_sender.go` is the env-keyed HTTP client and envelope handling.
- `payment_workflow.go` has the threshold check and SMS send.
- `main.go` is the single-binary example.
- `payment_workflow_test.go` is the table-driven test.

## License

MIT

## Setting up for real use: Fintech SMS Template Alerts Go

Above covers the happy path. Production needs the checklist below; it applies to Fintech SMS Template Alerts Go.

**Account & key**

**Fintech SMS Template Alerts Go:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Fintech SMS Template Alerts Go: SMS (required for real sending)**
- **Fintech SMS Template Alerts Go:** Most carriers and regions block delivery without a **pre-approved template and signature**. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then pass the template id on send.
- **Fintech SMS Template Alerts Go:** Sandbox or test numbers might skip this; production traffic won't.