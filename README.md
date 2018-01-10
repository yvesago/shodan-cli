# shodan-cli


Simple golang Shodan command line client with default query.


## Usage

To start working with Shodan you need an API key. You can do this at [https://www.shodan.io](https://www.shodan.io).

Use the API key in `$SHODAN` environment variable.


```bash
Usage of ./shodan-cli:
  -b	black & white, no color
  -c	compact, no detail
  -i string
    	ip [192.168.0.1]
  -n string
    	net [192.168.0.0/24]
  -q string
    	query ['!http']
```

On first call `shodan-cli` will ask an optional default query stored in `.shoddanrc`.


### Query sample
![Shodan Query](img/ShodanQuery.png)


### Network query sample
![Shodan Net Query](img/ShodanNetQuery.png)


## Build

TODO


## Links
* [Shodan.io](http://shodan.io)
* [Golang Shoddan API](http://github.com/ns3777k/go-shodan)


## Licence

MIT License

Copyright (c) 2018 Yves Agostini

<yves+github@yvesago.net>
