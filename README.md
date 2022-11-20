# Deformd

Simple and self-hosted form as a service.

> ⚠ Deformd is currently in a very alpha stage ! Expect breaking changes...

## Install

### Manually

Download the pre-compiled binaries from the [releases page](https://github.com/Bornholm/deformd/releases) and copy them to the desired location.

### Bash script

```
curl -sfL https://raw.githubusercontent.com/Bornholm/deformd/master/misc/script/install.sh | bash
```

It will download `deformd` to your current directory.

#### Install script environment variables

|Name|Description|Default|
|----|-----------|-------|
|`DEFORMD_VERSION`|Deformd version to download|`latest`|
|`DEFORMD_DESTDIR`|Deformd destination directory|`.`|

## Usage

```
deformd run -c config.yml
```

## Configuration

[View an example](internal/config/testdata/config.yml)

```yaml
                            # Logger configuration
logger: 
  level: 2                  # Logging level (0: DEBUG, 1: INFO, 2: WARN, 3: ERROR, 4: CRITICAL)
  format: human             # Logging format, "human" or "json"

http:                       # Web server configuration
  host: "0.0.0.0"           # Listening host
  port: 3000                # Listening port 

forms:                      # Forms configuration
  <FORM_NAME>: <FORM_SECTION>
```

### `<FORM_SECTION>`

```yaml
<FORM_NAME>:
  fields: <FIELDS_SECTION>
  handler: <HANDLER_SECTION>
```

### `<FIELDS_SECTION>`

```yaml
fields:
  - <TEXT_FIELD_SECTION> | <NUMBER_FIELD_SECTION> | <EMAIL_FIELD_SECTION> | <SUBMIT_FIELD_SECTION>
```

#### `<TEXT_FIELD_SECTION>`

> TODO

#### `<NUMBER_FIELD_SECTION>`

> TODO

#### `<EMAIL_FIELD_SECTION>`

> TODO
#### `<SUBMIT_FIELD_SECTION>`

> TODO

### `HANDLER_SECTION`

```yaml
handler:
  script: <HANDLER_SCRIPT_SECTION>
  config: <HANDLER_CONFIG_SECTION>
```

#### `<HANDLER_SCRIPT_SECTION>`

> TODO

#### `<HANDLER_CONFIG_SECTION>`

```yaml
config:
  modules: <HANDLER_CONFIG_MODULES_SECTION>
```

#### `<HANDLER_CONFIG_MODULES_SECTION>

```yaml
modules:
  email: <HANDLER_EMAIL_MODULE_SECTION>
  params: <HANDLER_PARAMS_MODULE_SECTION>
```

#### `<HANDLER_EMAIL_MODULE_SECTION>`

```yaml
email:
  host: stmp-host.org
  port: 21
  username: my-username
  password: my-password
  insecureSkipVerify: false
  useSSL: true
  authType: PLAIN
  tlsPolicy: 1
```

#### `<HANDLER_PARAMS_MODULE_SECTION>`

```yaml
params:
  <KEY>: <VALUE>
```

## Changelog

[See `CHANGELOG.md`](./CHANGELOG.md)

## Licence

AGPL-3.0
