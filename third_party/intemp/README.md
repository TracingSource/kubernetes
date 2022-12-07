# intemp

A bash script to execute a command within a temporary work directory.


## Dependencies

Requires: mktemp


## Install

```
git clone https://github.com/karlkfi/intemp
cd intemp
make install
```

or

```
curl -o- https://raw.githubusercontent.com/karlkfi/intemp/master/install.sh | bash
```

## Usage

```
intemp.sh [-t prefix] "<command>"
```

Example (install intemp using intemp):

```
intemp.sh -t intemp "git clone https://github.com/karlkfi/intemp . && make install"
```


## License

Copyright 2015 Karl Isenberg

