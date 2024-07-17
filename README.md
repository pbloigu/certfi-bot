# CERT-FI Alert Mastodon bot

Watches the [CERT-FI](https://www.kyberturvallisuuskeskus.fi/en/our-activities/cert) security alert feed and toots new alerts to a Mastodon instance.

## Description

Periodically (period configurable) reads the CERT-FI alert RSS feed (address configurable) and if new entries found after the last reading, constructs a Mastodon status message, which is then publised to a Mastodon instance (address and login details configurable).

This is a personal Golang learning project and as such likely not much of use for anyone else. Using this as an example of anything but maybe bad code is highly discouraged. 
## Getting Started


### Installing

* Download the latest release for your architecture
* Write a configuration file (in .yaml format) and place it in /etc/certfi-bot/config.yml. Make it root readable only. See the [test config file](./config.yml) for reference.

### Executing program
Execute the binary. It will remain running until killed and keep checking the CERT-FI rss feed at configured intervals.

## Help

I guess you can file an issue....


## Authors

Yours truly.

## License

This project is licensed under the CC0 1.0 Universal License - see the LICENSE.md file for details
