#!/bin/bash
protoc --go_out=. --go_opt=paths=source_relative game.proto error.proto enum.proto
protoc -I=. --cpp_out=D:\workspace\cocosspace\tgame\Classes game.proto error.proto enum.proto
