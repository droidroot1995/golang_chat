#!/bin/bash

cd "./Client"
go build
cd ".."
cd "./Server"
go build
cd ".."