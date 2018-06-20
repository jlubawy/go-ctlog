SRC_DIR := ./../../src

CC := gcc
CFLAGS += -s -g -O3 -Wall -I. -I$(SRC_DIR)

# https://stackoverflow.com/questions/6687630/how-to-remove-unused-c-c-symbols-with-gcc-and-ld
.PHONY: all
all:
	go install github.com/jlubawy/go-ctlog/...
	cmodule json . > cmodule_indices.json
	cmodule header cmodule_indices.json > cmodule_indices.h
	ctlog json cmodule_indices.json > ctlog_dict.json
	$(CC) $(CFLAGS) -o main $(SRC_DIR)/ctlog.c main.c
	$(CC) $(CFLAGS) -o main_ctlog $(SRC_DIR)/ctlog.c main_ctlog.c
