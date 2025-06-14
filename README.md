Mailhook
================

A simple mailhook server that listens for incoming emails and transform them into a webhook request.

## Configure

You can configure multiple webhook by using the following config example.
You can specify multiple from/to, and the action associated with these from/to.
you can add multiple actions for the same from/to combination.

```jsonc
[
  {
    "from": [
      "decima@zeus"
    ],
    "to": [
      "do@localhost"
    ],
    "action": {
      "url": "http://localhost:18000/00000000-0000-0000-0000-000000000000",
      "method": "POST",
      "format": "json",
      "body": {
        "subject": "{{subject}}",
        "from": "{{from}}",
        "to": "{{to}}",
        "text": "{{text}}",
        "some": "other",
        "json": "data"
      }
    }
  }
]
```


## How to use

With docker-compose: 
```yaml
services:
  mailhook:
    image: decima/mailhook:0.1.1
    ports:
      - "25:2525"
    volumes:
      - ./config.json:/app/config.json
```

## Development and local testing purpose
For this to run, you'll need:
- Docker (with compose, but let's be honest, if you don't have docker compose, it may be time to upgrade your stack)
- nix (optional, but recommended for development)
- swaks (a tool to send emails, useful for testing)

run
```shell
make local.start.deps
```
to start local deps. (not mandatory, but useful for development)
then you can go to [localhost:18000/s/00000000-0000-0000-0000-000000000000](http://localhost:18000/s/00000000-0000-0000-000000000000) to see the webhook endpoint.

Then you can run ```make local.run``` which will start the server on port 2525.

if you want to test you can use : 
```shell
make test.mail
make test.mail.full
```


### Using nix
you can run
```shell
make
```
and it will start deps and install swaks.
