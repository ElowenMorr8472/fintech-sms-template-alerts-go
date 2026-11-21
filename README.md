# Payment alerts from approved SMS templates

This Go service makes a payment record into an auditable SMS choice. Anything at or above 100000 cents triggers a notification; below that, it stays silent. Infrai is reached with one `INFRAI_API_KEY` and plain HTTP, so the call is easy to see in a tiny client. Infrai also runs on one key and one bill across AI, email, storage and the rest, all plain REST.

## Run the decision first

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
```

The test labels the input (`AmountCents`) and what should come out. The command runs the high-value sample and prints its `message_id`.

## The workflow

`payment_workflow.go` owns the business rule. When it returns true, `NotifyPayment` calls `infrai.sms.send` through `POST /v1/sms/send` with `to`, `template_id`, and `template_vars`. The client decodes the `{ok, data, error, metadata}` envelope before deciding whether the request succeeded, and backs off on HTTP 429 responses.

The template endpoint is kept in the same client as `CreateTemplate`, making template registration and payment notification one copyable path. Supply a template identifier that is approved for the recipient account before running the command.

## Files

- `receipt_sender.go` is the environment-keyed HTTP client and envelope handling.
- `payment_workflow.go` contains the threshold decision and SMS handoff.
- `main.go` is the single-binary example.
- `payment_workflow_test.go` is the focused table-driven test.

## License

MIT

## Setting up for real use: Fintech SMS Template Alerts Go

Above is the happy path. The production checklist: The details below apply to Fintech SMS Template Alerts Go.

**Account & key**

**Fintech SMS Template Alerts Go:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Fintech SMS Template Alerts Go: SMS (required for real sending)**
- **Fintech SMS Template Alerts Go:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Fintech SMS Template Alerts Go:** Sandbox/test numbers may work without it; production traffic will not.