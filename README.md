DoozerD Picker
==============
This picks the right namespace for a new or recovering Doozerd to bind to from N other servers that might or might not still exist.

```
go get github.com/musicglue/doozerd-picker
go install github.com/musicglue/doozerd-picker
doozerd-picker -servers="oneserver,twoserver" -port 8046 -timeout 500 -protocol tcp
```