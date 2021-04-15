# garmin-intl2cn

**UNSTABLE**

## Usage notice

Requesting unofficial APIs and simulated website login behaviors can be suspicious to Garmin.

Use it AT YOUR OWN RISK.

## Usage

```
cp config/config.sample config/config.go
# Filling config.go with garmin international and CN account info
cd main
go build -ldflags '-s -w' -o garmin-intl2cn

# Specific port to listen. default is 38080.
PORT=38080 ./garmin-intl2cn

# Sync latest activities(up to 3 activities) of garmin international account to CN account
# It will try to log in to Garmin website in ervery sync process since login session persistence is not implemented.
curl 'http://localhost:38080/api/sync'
```

## Thanks

[tapiriik](https://github.com/cpfair/tapiriik)

[garmin-connect](https://github.com/abrander/garmin-connect)

[garminexport](https://github.com/petergardfjall/garminexport)