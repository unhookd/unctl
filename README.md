# UnhookD

UnhookD is a release controller designed for safe integration with external deployment triggers.

## Usage
### Server commands
```
/usr/bin/unhookd zero-trust-server # Runs unhookd in secure zero-trust model mode
```

### Client commands

```
/usr/bin/unhookd deploy [deployment] [release] [sha] [flags] # deploys an application with the given args
```

## Deploying an application with Unhookd
Unhookd accepts deploy requests in zero trust mode, see [ZERO_TRUST_MODE.md](./ZERO_TRUST_MODE.md)  

### Zero Trust Mode
Zero trust mode should be used for environments such as production and staging.

See [ZERO_TRUST_MODE.md](./ZERO_TRUST_MODE.md) for more information on how to configure an application to deploy with zero trust.

## Developing
See [DEVELOPING.md](./DEVELOPING.md) for more information on how to develop on Unhookd.

