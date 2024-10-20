# BlackSniper 
<p align="center">
<img src="https://i.pinimg.com/736x/a3/e3/46/a3e3468de0f3789636dd1dab0fee558c.jpg", width="500", height="500">
</p>

A Discord selfbot for automatically redeeming Nitro codes.
Based on the old Crown sniper 
Using fasthttp + webhook support

# Requirements

- Dependencies
- A valid Discord token
- A valid webhook URL
- Go (version 1.14 or higher)


# Usage
The bot can be used in two ways:

- Run the bot with the configuration file config.toml
- Run the bot with the -t and -w parameters to specify the Discord token and webhook URL   
Example : 
```bash
go run main.go -t token -w gamingwebhook
