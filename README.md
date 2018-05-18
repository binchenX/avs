# AVS - create Android device easier and faster

[![GoDoc](https://godoc.org/github.com/pierrchen/avs?status.svg)](https://godoc.org/github.com/pierrchen/avs)


## The Problems 

There are two major problems nowadays when trying to building an Android Devices.

### 1. Lack of specification and comprehensive documentation.

Here is a dump-out of the consequence of that:

* Don't know where to start when adding a new device.
* Don't know which configurations are mandatory/minimal.
* Don't know exactly what a configuration is for and what it will impact.
* Not obvious which configuration should be put into which configure files and why.
* Not obvious for a certain feature, what configurations are available.
* The configurations are not logically organized.
* You can put configurations in many different places. No best practice.
* Duplicated configurations.
* Obsoleted configuration still in use for a new product.
* Hard learned knowledge can't be shared and become a one-time thing.

### 2. Long turnaround time for configurations change.

Depend on the type of change as well as the type of the error, it takes 10 minutes to 2 hours to detect a device configuration error, which could have been caught in seconds once we solve the 1st problem.

## An Attempted Solution  

To improve the situation, I came up with `AVS`, a short for `Andriod Vendor Build Specification`, for whatever it means.

It provides two things: 1) A *Spec*, that specifies all the configurations; 2) A *tool* facilitate you building the device using JSON configuration file.

### 1. Spec

The [spec](https://github.com/pierrchen/avs/tree/master/spec) is written in `go` and thanks to that fact, it comes with a [Doc](https://godoc.org/github.com/pierrchen/avs/spec#Spec) automatically.

* It is *the* place for all the configurations, 
* It is organized in a top-down tree structure
* It specifies and documents exactly what is available and what they are for.

This aims to the place we all collaborate and improve overtime.

### 2. Tool

#### 2.1 Load `json` config file

Json strikes a good balance of human readable and machine readable, and it has good integration with ide/editors and provide much better user experience than Makefile.

#### 2.2 Configration validation

Catch error in seconds not hours later. 

* Ensure mandatory configurations are there
* Ensure the configurations are valid
* Ensure no conflicted configurations
* Cross-check different configurations. (e.g new device node should have new SELinux policy)
* Plugin based

**The goal is once you pass the schema validation, you can pass most of `VTS`.**

### 2.3 Generate the .mk file

* It generates the *correct* .mk file that the build system would expect.

* It has the best practice built-in and will be enforced in the validation stage. So follow the [convention](https://en.wikipedia.org/wiki/Convention_over_configuration), not one of the five different ways.

* It supports config fragments. 

  For example, you can have one file called `hal.wifi.overlay` which contains all the configurations needed for wifi. Drop that in the device directory, and do `avs update`, it will either add (when there is no wifi config in config.json) or override (when there is a wifi in the config.json) the wifi configration. It is much manageable than having multiple `mk` files and struggle with where to put what.

## Download

### Binary

[v1.0](https://github.com/pierrchen/avs/releases/download/v1.0/avs)

### Go get

```bash
$go get github.com/pierrchen/avs
```
To get started:
```
  $cd ${ANDROID_BUILD_TOP}/device
  $avs init  --vendor linaro --device awesome
  $cd linaro/awesome
  # modify the default config.json
  $avs udpate
```

see `avs -h`.

## AVS In Action

[Poplar Device](https://github.com/96boards-poplar/poplar-device) is created from [this configs](https://github.com/pierrchen/avs/tree/master/devices/poplar).
