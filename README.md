# Slowport-Restarter

Does your Speedport W 724V router also get slow if you do not pull the plug regularly? Are you also annoyed that there is no periodic restart capability built into it's management interface? If yes, you are [not](https://telekomhilft.telekom.de/t5/Telefonie-Internet/Speedport-W724V-alle-2-tage-neustart-notwendig-ausfall-telefon/td-p/1784426) [alone](https://telekomhilft.telekom.de/t5/Telefonie-Internet/Router-Speedport-W724v-automatisch-in-24-Stunden-neustarten/td-p/1287775)...

Here comes your solution. This repo contains a simple command line utility to automatically restart the router. Call it with your router's IP and password as arguments and your router will perform a reboot. Configure it as a cronjob and you are good to go.

## How to use

Locate and download a suitable binary in `go/build/slowport-restarter-<OS>-<Architecture>`. Then run it as follows (example for Mac OS):

```bash
slowport-restarter-darwin-amd64 --host speedport.ip --password <your router password>
```

If this does not work, use your router's IP rather than `speedport.ip` (may be necessary if you messed with your DNS settings).
