# UnhookD

UnhookD is a release controller designed for safe integration with external deployment triggers.

## Usage
### Server commands
```
/usr/bin/unhookd instastage-server # Runs unhookd in insecure legacy instastage mode
/usr/bin/unhookd zero-trust-server # Runs unhookd in secure zero-trust model mode
```

### Client commands

```
/usr/bin/unhookd deploy [deployment] [release] [sha] [flags] # deploys an application with the given args
```

## Deploying an application with Unhookd
Unhookd accepts deploy requests in two modes: zero trust and instastage.

### Zero Trust Mode
In instastage mode, Unhookd will relay a request to deploy a given `$project` & `$release` at a given `$sha`. Zero trust mode should be used for environments such as production and staging.

See [ZERO_TRUST_MODE.md](./ZERO_TRUST_MODE.md) for more information on how to configure an application to deploy with zero trust.

## Developing
See [DEVELOPING.md](./DEVELOPING.md) for more information on how to develop on Unhookd.

