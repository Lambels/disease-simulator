[![Version](https://img.shields.io/badge/goversion-1.18.x-blue.svg)](https://golang.org)
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>

# Disease-Simulator

Simulate a configurable disease on a configurable sample size.

## Build

Mac/Linux:

```
go build -o dis-sim *.go
```

Windows:

```
go build -o dis-sim.exe .
```

For particular platform:

```
env GOOS=linux GOARCH=amd64 go build -o dis-sim *.go
```

## Run

```
./dis-sim \
-population-size=10
-population-infected=5
-population-max-interactions=3
-infection-rate=0.7 # %70
-process-generators=1
-process-processors=1
-process-processor-timeout=10 # 10 seconds
```

## All Flags

```
./dis-sim -help
Usage of dis-sim:
-infection-rate float
    infection-rate determines the chance of disease transmition, ie: 0.5, 0.25, 1.0 (default 0.5)
-population-infected int
    population-infected determines the number of starting infected sample (default 5)
-population-max-interactions int
    population-max-interactions determines the max number of interactions a pacient can have (default 3)
-population-size int
    population-size determines the sample size (default 10)
-process-generators int
    process-generators determines the number of generators used (default 1)
-process-processor-timeout int
    process-processor-timeout determines the number of seconds used as processor timeout (default 10)
-process-processors int
    process-processors determines the number of processors used (default 1)
```