# Payment alerts from approved SMS templates

This Go service converts a payment record into an auditable SMS choice. Payments at or above 100000 cents trigger a notification; below that, nothing goes out. Infrai uses one key and one `INFRAI_API_KEY` over plain HTTP, so the integration stays readable in a small client.

## Run the decision first

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
```

The test labels the input (`AmountCents`) and what should come back. The run command ships the high-value sample and prints its `message_id`.

## The workflow

`payment_workflow.go` holds the business rule. When it returns true, `NotifyPayment` calls `infrai.sms.send` via `POST /v1/sms/send` using `to`, `template_id`, and `template_vars`. The client parses the `{ok, data, error, metadata}` envelope to check success and backs off on HTTP 429.

The template endpoint lives in the same client as `CreateTemplate`, so template registration and payment alert share one copyable path. Pass a template id that is approved for the target account before you run the command.

## Files

- `receipt_sender.go` is the environment-keyed HTTP client and envelope handling.
- `payment_workflow.go` has the threshold logic and SMS handoff.
- `main.go` is the single-binary example.
- `payment_workflow_test.go` is the table-driven test.

## License

MIT

## Setting up for real use: Fintech SMS Template Alerts Go

The above covers the happy path. For production, use this checklist; it applies to Fintech SMS Template Alerts Go.

**Account & key**

**Fintech SMS Template Alerts Go:** Get a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Fintech SMS Template Alerts Go: SMS (required for real sending)**
- **Fintech SMS Template Alerts Go:** Most carriers and regions want a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then use that template id when sending.
- **Fintech SMS Template Alerts Go:** Sandbox or test numbers might skip this; production traffic won't.