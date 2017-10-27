# lunch-buddies
Creating randomized groups to get people to meet each other across the company.

## Prerequisites
Download the newest go lang library [here](https://golang.org/doc/install?download=go1.5.windows-amd64.msi2). Run `go get github.com/golang/glog`

## How to Run
Create a csv file. Currently this file has to have a pretty specific format. Each column should have a header. In order, it goes: `Name, Team, Email`. If you do not add the headers, the file will not process, and you will be shown an error.

In a terminal navigate to where the project is located, and run `go run buddy.go --input_file <path to csv file>`. Other configurable options are as follows:
* __--group_size <integer>:__ Defaults to 6 if not set. When sorting the inputed people, will sort them into groups of this size.
* __--randomize <boolean>:__ Defaults to true. When set to false, will group people in the given order.
* __--input_file-- <path to csv file>:__ Defaults to `./tmp/workers.csv`. You can point this to any csv file.
