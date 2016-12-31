# sproutely
Find and submit coupons for Sprouts supermarket sans app

Saves needing to manually put 100's of coupons on to one's account using the Sprouts app.

You'll need to have an account with Sprouts, and then ideally download the phone app so you can actually get value from this.

usage:

```
  [crankyflamingo@wintermute sproutely]$ ./sproutely 
    -login
        Used to regenerate token, by logging in with username and password. Tokens are typically valid for months
    -update
        Will log into site, gather coupons, and apply to account
```

config.json:

```
  {"User":"crankyflamingo@github","Pass":"Passw0rd!","Token":""}
```

Token field will be filled in upon first login. 

